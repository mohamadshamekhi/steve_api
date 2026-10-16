package partnerservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	partnerparam "github.com/abolfazlnorzad/steve/param/partner"
	paymentparam "github.com/abolfazlnorzad/steve/param/payment"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// --- shop ports ---

// CartProvider is the partner cart, keyed by the partner's backoffice user id.
type CartProvider interface {
	GetUserActiveCart(ctx context.Context, userID string) (*entity.Cart, error)
	UpsertItemWithQuantity(ctx context.Context, userID string, item entity.CartItem, quantity int) error
	RemoveFromCart(ctx context.Context, userID string, productID string, variantID string) error
	UpdateStatus(ctx context.Context, cartID string, status entity.CartStatus) error
}

// OrderWriter persists partner orders. Satisfied by *mongodborder.DB.
type OrderWriter interface {
	Create(ctx context.Context, order entity.Order) error
}

// ProductStockProvider decrements variant stock after an order is registered.
// Satisfied by productservice.BackofficeService.
type ProductStockProvider interface {
	UpdateProductVariantQuantityAfterRegisterOrder(ctx context.Context, req productparam.UpdateProductVariantQuantityAfterRegisterOrderRequest) (productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse, error)
}

// PaymentProvider drives online payment for partner orders. Satisfied by
// paymentservice.Service.
type PaymentProvider interface {
	GoToGateway(ctx context.Context, req paymentparam.GoToGatewayRequest) (paymentparam.GoToGatewayResponse, error)
	VerifyPayment(ctx context.Context, req paymentparam.VerifyPaymentRequest) (paymentparam.VerifyPaymentResponse, error)
}

// ShopService is the partner-panel-facing service. Its caller is the partner's
// restricted BACKOFFICE account, so every method resolves the partner record by
// its backoffice user id.
type ShopService struct {
	repo     Repository
	products ProductProvider
	orders   PartnerOrderProvider
	cart     CartProvider
	orderW   OrderWriter
	stock    ProductStockProvider
	payment  PaymentProvider
	resolver PartnerPriceResolver
	cfg      Config
}

func NewShopService(
	repo Repository,
	products ProductProvider,
	orders PartnerOrderProvider,
	cart CartProvider,
	orderW OrderWriter,
	stock ProductStockProvider,
	payment PaymentProvider,
	resolver PartnerPriceResolver,
	cfg Config,
) ShopService {
	return ShopService{
		repo:     repo,
		products: products,
		orders:   orders,
		cart:     cart,
		orderW:   orderW,
		stock:    stock,
		payment:  payment,
		resolver: resolver,
		cfg:      cfg,
	}
}

func (s ShopService) requirePartner(ctx context.Context, op, backofficeUserID string) (entity.Partner, error) {
	if !s.cfg.Enabled {
		return entity.Partner{}, richerror.New(op).WithKind(richerror.KindForbidden).WithMessageKey("partner.program_disabled")
	}
	p, err := s.repo.FindByBackofficeUserID(ctx, backofficeUserID)
	if err != nil {
		return entity.Partner{}, richerror.New(op).WithWarpError(err)
	}
	if !p.IsApproved() {
		return entity.Partner{}, richerror.New(op).WithKind(richerror.KindForbidden).WithMessageKey("partner.not_approved")
	}
	return p, nil
}

// productEligible reports whether the partner discount applies to a product.
// An explicit admin exclusion always wins; otherwise the product is eligible
// if opted in, or if the program runs in all-products mode.
func (s ShopService) productEligible(p entity.Product) bool {
	if p.PartnerExcluded {
		return false
	}
	return p.PartnerEligible || s.cfg.AllProductsEligible
}

// priceVariant resolves the partner price of one product variant.
func (s ShopService) priceVariant(p entity.Product, variantID string, partnerPercent uint) (entity.Variant, PriceResult, bool) {
	for _, v := range p.Variants {
		if v.ID == variantID {
			return v, s.resolver.Resolve(PriceInput{
				BasePrice:              v.Price,
				ProductEligible:        s.productEligible(p),
				ProductDiscountPercent: p.PartnerDiscountPercent,
				PartnerDiscountPercent: partnerPercent,
			}), true
		}
	}
	return entity.Variant{}, PriceResult{}, false
}

// Me returns the partner's own status and discount for the panel header.
func (s ShopService) Me(ctx context.Context, backofficeUserID string) (partnerparam.MeResponse, error) {
	const op = "partnerservice.ShopService.Me"
	p, err := s.repo.FindByBackofficeUserID(ctx, backofficeUserID)
	if err != nil {
		if richerror.GetKind(err) == richerror.KindNotFound {
			return partnerparam.MeResponse{HasPartnership: false}, nil
		}
		return partnerparam.MeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return partnerparam.MeResponse{
		HasPartnership:  true,
		Status:          string(p.Status),
		DiscountPercent: p.DiscountPercent,
		CompanyName:     p.CompanyName,
	}, nil
}

// Catalog lists published products annotated with the partner price for every
// eligible product, using the shared PartnerPriceResolver.
func (s ShopService) Catalog(ctx context.Context, backofficeUserID string, req partnerparam.CatalogListRequest) (partnerparam.CatalogListResponse, error) {
	const op = "partnerservice.ShopService.Catalog"
	partner, err := s.requirePartner(ctx, op, backofficeUserID)
	if err != nil {
		return partnerparam.CatalogListResponse{}, err
	}

	filter := req.Filter
	filter.Filter = append(filter.Filter, map[string]any{"is_published": "true"})

	list, err := s.products.All(ctx, filter, req.Pagination)
	if err != nil {
		return partnerparam.CatalogListResponse{}, richerror.New(op).WithWarpError(err)
	}

	products := make([]partnerparam.CatalogProduct, 0, len(list.Products))
	for _, p := range list.Products {
		cp := partnerparam.CatalogProduct{
			ID:              p.ID,
			Name:            p.Name,
			Slug:            p.Slug,
			URL:             p.GetMainGallery().URL,
			PartnerEligible: s.productEligible(p),
			Variants:        make([]partnerparam.CatalogVariant, 0, len(p.Variants)),
		}
		for _, v := range p.Variants {
			res := s.resolver.Resolve(PriceInput{
				BasePrice:              v.Price,
				ProductEligible:        s.productEligible(p),
				ProductDiscountPercent: p.PartnerDiscountPercent,
				PartnerDiscountPercent: partner.DiscountPercent,
			})
			cp.AppliedPercent = res.AppliedPercent
			cp.Variants = append(cp.Variants, partnerparam.CatalogVariant{
				ID:             v.ID,
				Price:          v.Price,
				PartnerPrice:   res.FinalPrice,
				DiscountAmount: res.DiscountAmount,
				Quantity:       v.Quantity,
				Warranty:       v.Warranty,
			})
		}
		products = append(products, cp)
	}

	return partnerparam.CatalogListResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(list.TotalCount),
		},
		Products: products,
	}, nil
}

// AddToCart puts `quantity` units of a variant in the partner's cart at the
// partner price (absolute set, not increment).
func (s ShopService) AddToCart(ctx context.Context, backofficeUserID string, req partnerparam.AddToCartRequest) (partnerparam.CartResponse, error) {
	const op = "partnerservice.ShopService.AddToCart"
	partner, err := s.requirePartner(ctx, op, backofficeUserID)
	if err != nil {
		return partnerparam.CartResponse{}, err
	}
	qty := req.Quantity
	if qty < 1 {
		qty = 1
	}

	p, err := s.products.FindByID(ctx, req.ProductID)
	if err != nil {
		return partnerparam.CartResponse{}, richerror.New(op).WithWarpError(err)
	}
	variant, priced, ok := s.priceVariant(p, req.VariantID, partner.DiscountPercent)
	if !ok {
		return partnerparam.CartResponse{}, richerror.New(op).
			WithKind(richerror.KindNotFound).WithMessageKey("partner.variant_not_found")
	}
	if variant.Quantity < uint64(qty) {
		return partnerparam.CartResponse{}, richerror.New(op).
			WithKind(richerror.KindBadRequest).WithMessageKey("partner.insufficient_stock")
	}

	item := entity.CartItem{
		ProductID:         p.ID,
		VariantID:         variant.ID,
		ProductName:       p.Name,
		VariantAttributes: variant.VariantAttribute,
		PriceAtTime:       int64(priced.FinalPrice),
	}
	if err := s.cart.UpsertItemWithQuantity(ctx, backofficeUserID, item, qty); err != nil {
		return partnerparam.CartResponse{}, richerror.New(op).WithWarpError(err)
	}
	return s.GetCart(ctx, backofficeUserID)
}

// RemoveFromCart drops a line from the partner's cart.
func (s ShopService) RemoveFromCart(ctx context.Context, backofficeUserID string, req partnerparam.RemoveFromCartRequest) (partnerparam.CartResponse, error) {
	const op = "partnerservice.ShopService.RemoveFromCart"
	if _, err := s.requirePartner(ctx, op, backofficeUserID); err != nil {
		return partnerparam.CartResponse{}, err
	}
	if err := s.cart.RemoveFromCart(ctx, backofficeUserID, req.ProductID, req.VariantID); err != nil {
		return partnerparam.CartResponse{}, richerror.New(op).WithWarpError(err)
	}
	return s.GetCart(ctx, backofficeUserID)
}

// GetCart returns the partner's active cart, re-priced against the current
// partner discount so stale prices never leak into checkout.
func (s ShopService) GetCart(ctx context.Context, backofficeUserID string) (partnerparam.CartResponse, error) {
	const op = "partnerservice.ShopService.GetCart"
	partner, err := s.requirePartner(ctx, op, backofficeUserID)
	if err != nil {
		return partnerparam.CartResponse{}, err
	}

	cart, err := s.cart.GetUserActiveCart(ctx, backofficeUserID)
	if err != nil {
		if richerror.GetKind(err) == richerror.KindNotFound {
			return partnerparam.CartResponse{Items: []partnerparam.CartItemView{}}, nil
		}
		return partnerparam.CartResponse{}, richerror.New(op).WithWarpError(err)
	}

	resp := partnerparam.CartResponse{CartID: cart.ID, Items: make([]partnerparam.CartItemView, 0, len(cart.Items))}
	for _, it := range cart.Items {
		p, pErr := s.products.FindByID(ctx, it.ProductID)
		if pErr != nil {
			continue
		}
		variant, priced, ok := s.priceVariant(p, it.VariantID, partner.DiscountPercent)
		if !ok {
			continue
		}
		line := priced.FinalPrice * uint64(it.Quantity)
		resp.Items = append(resp.Items, partnerparam.CartItemView{
			ProductID:   p.ID,
			VariantID:   variant.ID,
			ProductName: p.Name,
			URL:         p.GetMainGallery().URL,
			Quantity:    it.Quantity,
			BasePrice:   variant.Price,
			UnitPrice:   priced.FinalPrice,
			LineTotal:   line,
			InStock:     variant.Quantity >= uint64(it.Quantity),
		})
		resp.Subtotal += line
		resp.SavedTotal += priced.DiscountAmount * uint64(it.Quantity)
	}
	return resp, nil
}

// Checkout turns the partner cart into an order. `in_person` leaves it pending
// for the counter; `online` also opens the payment gateway and returns its URL.
func (s ShopService) Checkout(ctx context.Context, backofficeUserID string, req partnerparam.CheckoutRequest) (partnerparam.CheckoutResponse, error) {
	const op = "partnerservice.ShopService.Checkout"
	partner, err := s.requirePartner(ctx, op, backofficeUserID)
	if err != nil {
		return partnerparam.CheckoutResponse{}, err
	}

	method := req.PaymentMethod
	if method != partnerparam.PaymentMethodOnline && method != partnerparam.PaymentMethodInPerson {
		return partnerparam.CheckoutResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).WithMessageKey("partner.invalid_payment_method")
	}

	cart, err := s.cart.GetUserActiveCart(ctx, backofficeUserID)
	if err != nil {
		if richerror.GetKind(err) == richerror.KindNotFound {
			return partnerparam.CheckoutResponse{}, richerror.New(op).WithKind(richerror.KindBadRequest).WithMessageKey("partner.cart_empty")
		}
		return partnerparam.CheckoutResponse{}, richerror.New(op).WithWarpError(err)
	}
	if len(cart.Items) == 0 {
		return partnerparam.CheckoutResponse{}, richerror.New(op).WithKind(richerror.KindBadRequest).WithMessageKey("partner.cart_empty")
	}

	var subtotal, saved uint64
	items := make([]entity.OrderItem, 0, len(cart.Items))
	for _, it := range cart.Items {
		p, pErr := s.products.FindByID(ctx, it.ProductID)
		if pErr != nil {
			return partnerparam.CheckoutResponse{}, richerror.New(op).WithWarpError(pErr)
		}
		variant, priced, ok := s.priceVariant(p, it.VariantID, partner.DiscountPercent)
		if !ok {
			return partnerparam.CheckoutResponse{}, richerror.New(op).
				WithKind(richerror.KindNotFound).WithMessageKey("partner.variant_not_found")
		}
		if variant.Quantity < uint64(it.Quantity) {
			return partnerparam.CheckoutResponse{}, richerror.New(op).
				WithKind(richerror.KindBadRequest).WithMessageKey("partner.insufficient_stock")
		}
		subtotal += priced.FinalPrice * uint64(it.Quantity)
		saved += priced.DiscountAmount * uint64(it.Quantity)
		items = append(items, entity.OrderItem{
			ProductID:         p.ID,
			VariantID:         variant.ID,
			Quantity:          uint(it.Quantity),
			ProductName:       p.Name,
			VariantAttributes: variant.VariantAttribute,
			PriceAtTime:       uint(priced.FinalPrice),
		})
	}

	order := entity.Order{
		UserID:         backofficeUserID,
		CartID:         cart.ID,
		Items:          items,
		SubtotalAmount: uint(subtotal),
		FinalAmount:    uint(subtotal),
		Status:         entity.OrderStatusPendingPayment,
		IsPartnerOrder: true,
		PartnerID:      partner.ID,
	}
	if method == partnerparam.PaymentMethodInPerson {
		order.PurchaseType = entity.PurchaseTypeInPerson
	} else {
		order.PurchaseType = entity.PurchaseTypeOnline
	}
	if saved > 0 {
		order.AppliedDiscounts = []entity.AppliedDiscount{{
			DiscountCode: "partner_price",
			Amount:       int64(saved),
		}}
	}
	order.BeforeCreate()

	if err := s.orderW.Create(ctx, order); err != nil {
		return partnerparam.CheckoutResponse{}, richerror.New(op).WithWarpError(err)
	}
	if _, err := s.stock.UpdateProductVariantQuantityAfterRegisterOrder(ctx, productparam.UpdateProductVariantQuantityAfterRegisterOrderRequest{Order: order}); err != nil {
		return partnerparam.CheckoutResponse{}, richerror.New(op).WithWarpError(err)
	}
	if err := s.cart.UpdateStatus(ctx, cart.ID, entity.PendingPayment); err != nil {
		return partnerparam.CheckoutResponse{}, richerror.New(op).WithWarpError(err)
	}

	resp := partnerparam.CheckoutResponse{OrderID: order.ID, PaymentMethod: method}
	if method == partnerparam.PaymentMethodOnline {
		gw, gErr := s.payment.GoToGateway(ctx, paymentparam.GoToGatewayRequest{OrderID: order.ID, UserID: backofficeUserID})
		if gErr != nil {
			return partnerparam.CheckoutResponse{}, richerror.New(op).WithWarpError(gErr)
		}
		resp.PaymentURL = gw.PaymentURL
		resp.GatewayTransactionID = gw.GatewayTransactionID
	}
	return resp, nil
}

// VerifyPayment finalises an online partner payment via the gateway callback.
func (s ShopService) VerifyPayment(ctx context.Context, backofficeUserID string, req partnerparam.VerifyPaymentRequest) (paymentparam.VerifyPaymentResponse, error) {
	const op = "partnerservice.ShopService.VerifyPayment"
	if _, err := s.requirePartner(ctx, op, backofficeUserID); err != nil {
		return paymentparam.VerifyPaymentResponse{}, err
	}
	return s.payment.VerifyPayment(ctx, paymentparam.VerifyPaymentRequest{GatewayTransactionID: req.GatewayTransactionID})
}

// MyOrders is the partner's own purchase report.
func (s ShopService) MyOrders(ctx context.Context, backofficeUserID string, page commonparam.PaginationRequest) (partnerparam.MyOrdersResponse, error) {
	const op = "partnerservice.ShopService.MyOrders"
	if _, err := s.requirePartner(ctx, op, backofficeUserID); err != nil {
		return partnerparam.MyOrdersResponse{}, err
	}
	res, err := s.orders.GetUserOrderWithDetails(ctx, backofficeUserID, page)
	if err != nil {
		return partnerparam.MyOrdersResponse{}, richerror.New(op).WithWarpError(err)
	}
	return partnerparam.MyOrdersResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   page.GetPageSize(),
			PageNumber: page.GetPageNumber(),
			Total:      uint(res.TotalCount),
		},
		Orders: res.Orders,
	}, nil
}

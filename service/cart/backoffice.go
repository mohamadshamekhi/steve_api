package cartservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	orderservice "github.com/abolfazlnorzad/steve/service/order"
	"golang.org/x/sync/errgroup"
)

type BackofficeService struct {
	repo            CartBackofficeReadWriter
	uRepo           URepository
	pRepo           PRepository
	orderProvider   orderservice.BackofficeService
	productProvider ProductProvider
}

func NewBackofficeService(repo CartBackofficeReadWriter, uRepo URepository, pRepo PRepository, productProvider ProductProvider) BackofficeService {
	return BackofficeService{repo: repo, uRepo: uRepo, pRepo: pRepo, productProvider: productProvider}
}

func (s *BackofficeService) SetOrderProvider(orderProvider orderservice.BackofficeService) {
	s.orderProvider = orderProvider
}

func (s *BackofficeService) GetAll(ctx context.Context, req cartparam.GetAllCartsRequest) (cartparam.GetAllCartsResponse, error) {
	const op = "cartservice.GetAll"

	resp, err := s.repo.GetAll(ctx, req.Filter, req.Pagination)
	if err != nil {
		return cartparam.GetAllCartsResponse{}, richerror.New(op).WithWarpError(err)
	}

	if resp.TotalCount == 0 {
		return cartparam.GetAllCartsResponse{
			Pagination: commonparam.PaginationResponse{
				PageSize:   req.Pagination.GetPageSize(),
				PageNumber: req.Pagination.GetPageNumber(),
				Total:      0,
			},
			Carts: nil,
		}, nil
	}
	userIDs := make([]string, 0, len(resp.Carts))
	for _, cart := range resp.Carts {
		if cart.UserID != "" {
			userIDs = append(userIDs, cart.UserID)
		}
	}

	users, err := s.uRepo.FindManyByIDs(ctx, userIDs)
	if err != nil {
		logger.L().Error("cann't get users , users ids: %v - and errors : %v", userIDs, err)
	}

	var list []cartparam.BackofficeCartList

	for _, cart := range resp.Carts {
		u := entity.User{}
		for _, user := range users {
			if user.ID == cart.UserID {
				u = user
			}
		}
		list = append(list, cartparam.BackofficeCartList{
			ID:              cart.ID,
			Status:          cart.Status,
			UserID:          cart.UserID,
			UserPhoneNumber: u.PhoneNumber,
			UserName:        u.Name,
		})
	}

	return cartparam.GetAllCartsResponse{
		Carts: list,
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{PageSize: req.Pagination.PageSize, PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)},
	}, nil
}

func (s *BackofficeService) GetSingle(ctx context.Context, req cartparam.GetCartRequest) (cartparam.GetCartResponse, error) {
	const op = "cartservice.GetSingle"
	cart, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return cartparam.GetCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	ids := []string{cart.UserID}
	user, err := s.uRepo.FindManyByIDs(ctx, ids)
	if err != nil {
		return cartparam.GetCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	var cartItems []cartparam.CartItemDetail
	for _, item := range cart.Items {
		detail := cartparam.CartItemDetail{
			VariantID:         item.VariantID,
			Quantity:          item.Quantity,
			ProductName:       item.ProductName,
			VariantAttributes: item.VariantAttributes,
			PriceAtTime:       item.PriceAtTime,
			IsCustom:          item.IsCustom,
			Kind:              item.GetKind(),
			ManualDiscount:    item.ManualDiscount,
		}
		// A hand-typed line has no catalog product to look up, and asking for
		// one would fail the whole cart.
		if !item.IsCustom {
			p, err := s.pRepo.FindByID(ctx, item.ProductID)
			if err != nil {
				return cartparam.GetCartResponse{}, richerror.New(op).WithWarpError(err)
			}
			detail.ProductID = p.ID
			detail.Url = p.GetMainGallery().URL
		}
		cartItems = append(cartItems, detail)
	}

	c := cartparam.CartDetail{
		ID:     cart.ID,
		Status: string(cart.Status),
		User: struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
		}{
			ID:          cart.UserID,
			Name:        user[0].Name,
			PhoneNumber: user[0].PhoneNumber,
		},
		Items:     cartItems,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
	return cartparam.GetCartResponse{
		Cart: c,
	}, nil
}

func (s *BackofficeService) AddItemToCartBackoffice(ctx context.Context, req cartparam.AddItemToCartBackOfficeRequest) (cartparam.AddItemToCartBackOfficeResponse, error) {
	const op = "cartservice.AddItemToCartBackoffice"
	p, pErr := s.pRepo.FindByID(ctx, req.ProductID)
	if pErr != nil {
		return cartparam.AddItemToCartBackOfficeResponse{}, richerror.New(op).WithWarpError(pErr)
	}
	var price uint64
	var variants []entity.VariantAttribute
	for _, variant := range p.Variants {
		if variant.ID == req.VariantID {
			price = variant.Price
			variants = variant.VariantAttribute
		}
	}
	item := entity.CartItem{
		ProductID:         req.ProductID,
		VariantID:         req.VariantID,
		ProductName:       p.Name,
		VariantAttributes: variants,
		PriceAtTime:       int64(price),
	}
	if err := s.repo.AddToCart(ctx, req.UserID, item); err != nil {
		return cartparam.AddItemToCartBackOfficeResponse{}, richerror.New(op).WithWarpError(err)
	}

	return cartparam.AddItemToCartBackOfficeResponse{}, nil
}

func (s *BackofficeService) RemoveFromCartBackoffice(ctx context.Context, req cartparam.RemoveFromCartBackofficeRequest) (cartparam.RemoveFromCartBackofficeResponse, error) {
	const op = "cartservice.RemoveFromCartBackoffice"
	if err := s.repo.RemoveFromCart(ctx, req.UserID, req.ProductID, req.VariantID); err != nil {
		return cartparam.RemoveFromCartBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}

	return cartparam.RemoveFromCartBackofficeResponse{}, nil
}

func (s *BackofficeService) GetMyActiveCartBackoffice(ctx context.Context, req cartparam.MyActiveCartBackofficeRequest) (cartparam.MyActiveCartBackofficeResponse, error) {
	const op = "cartservice.GetMyActiveCartBackoffice"
	cart, err := s.repo.GetUserActiveCart(ctx, req.UserID)
	if err != nil {
		if richerror.GetKind(err) == richerror.KindNotFound {
			return cartparam.MyActiveCartBackofficeResponse{}, nil
		}
		return cartparam.MyActiveCartBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}

	var cartItems []cartparam.BackofficeCartItemDetail
	var itemsManualDiscount int64
	for _, item := range cart.Items {
		adj := item.ManualDiscount
		// A pure discount cannot exceed the line total; cap it.
		if adj > item.LineTotal() {
			adj = item.LineTotal()
		}
		detail := cartparam.BackofficeCartItemDetail{
			VariantID:         item.VariantID,
			Quantity:          item.Quantity,
			ProductName:       item.ProductName,
			VariantAttributes: item.VariantAttributes,
			PriceAtTime:       item.PriceAtTime,
			IsCustom:          item.IsCustom,
			Kind:              item.GetKind(),
			UnitCost:          item.UnitCost,
			ManualDiscount:    adj,
		}
		// A hand-typed line has no catalog product behind it.
		if !item.IsCustom {
			p, pErr := s.pRepo.FindByID(ctx, item.ProductID)
			if pErr != nil {
				return cartparam.MyActiveCartBackofficeResponse{}, richerror.New(op).WithWarpError(pErr)
			}
			detail.ProductID = p.ID
			detail.Url = p.GetMainGallery().URL
		}
		itemsManualDiscount += adj
		cartItems = append(cartItems, detail)
	}

	subtotal := cart.CalculateSubtotal()
	// payable is computed in int64 to allow markup to increase it.
	payable := int64(subtotal) - itemsManualDiscount
	cartManualDiscount := cart.ManualDiscount
	if int64(cartManualDiscount) > payable {
		cartManualDiscount = uint(payable)
	}
	payable -= int64(cartManualDiscount)

	c := cartparam.BackofficeCartDetail{
		ID:                  cart.ID,
		Status:              string(cart.Status),
		UserID:              cart.UserID,
		Items:               cartItems,
		CreatedAt:           cart.CreatedAt,
		UpdatedAt:           time.Now(),
		Subtotal:            subtotal,
		ManualDiscount:      cartManualDiscount,
		ItemsManualDiscount: itemsManualDiscount,
		PayableAmount:       payable,
	}
	return cartparam.MyActiveCartBackofficeResponse{
		Cart: c,
	}, nil
}

func (s *BackofficeService) DeleteItemFromCartBackoffice(ctx context.Context, req cartparam.DeleteItemFromCartBackofficeRequest) (cartparam.DeleteItemFromCartBackofficeResponse, error) {
	const op = "cartservice.DeleteItemFromCart"

	if err := s.repo.DeleteItemFromCart(ctx, req.UserID, req.VariantID); err != nil {
		return cartparam.DeleteItemFromCartBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return cartparam.DeleteItemFromCartBackofficeResponse{}, nil

}

func (s *BackofficeService) CreateCartBackoffice(ctx context.Context, req cartparam.CreateCartRequest) (cartparam.CreateCartResponse, error) {
	const op = "cartservice.CreateCartBackoffice"
	cart := entity.Cart{
		UserID: req.UserID,
	}
	cart.BeforeCreate()

	err := s.repo.CreateCart(ctx, cart)
	if err != nil {
		return cartparam.CreateCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	return cartparam.CreateCartResponse{}, nil
}

func (s *BackofficeService) UpdateStatus(ctx context.Context, req cartparam.UpdateCartStatusRequest) (cartparam.UpdateCartStatusResponse, error) {
	const op = "cartservice.UpdateStatus"

	cart, err := s.repo.FindByID(ctx, req.CartID)
	if err != nil {
		return cartparam.UpdateCartStatusResponse{}, richerror.New(op).WithWarpError(err)
	}

	if !entity.CanTransitionCartStatus(cart.Status, req.Status) {
		return cartparam.UpdateCartStatusResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("cart.invalid_status_transition").
			WithMeta(map[string]any{
				"from": cart.Status,
				"to":   req.Status,
			})
	}
	if uErr := s.repo.UpdateStatus(ctx, req.CartID, req.Status); uErr != nil {
		return cartparam.UpdateCartStatusResponse{}, richerror.New(op).WithWarpError(err)
	}
	return cartparam.UpdateCartStatusResponse{}, nil
}

func (s BackofficeService) CancelExpiredPendingCart(ctx context.Context, req cartparam.CancelExpiredPendingCartRequest) (cartparam.CancelExpiredPendingCartResponse, error) {
	const op = "orderservice.CancelExpiredPendingOrders"

	expiredBefore := time.Now().Add(-req.CanselExpiredPendingCart)

	carts, err := s.repo.GetExpiredPendingCart(ctx, expiredBefore)
	if err != nil {
		return cartparam.CancelExpiredPendingCartResponse{}, richerror.New(op).WithWarpError(err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	for _, cart := range carts {
		cart := cart
		cartID := cart.ID
		g.Go(func() error {
			var cErr error

			// Only stock-tracked variants can have stock returned to them;
			// an on-demand item or a service never held any, so crediting it
			// back would invent inventory that does not exist.
			variantIDs := make([]string, 0, len(cart.Items))
			for _, item := range cart.Items {
				// A hand-typed line never held catalog stock.
				if item.IsCustom {
					continue
				}
				variantIDs = append(variantIDs, item.VariantID)
			}
			variantsResp, vErr := s.productProvider.GetVariantsByIDs(gCtx, productparam.GetVariantsByIDsRequest{IDs: variantIDs})
			if vErr != nil {
				logger.L().Error("failed to resolve variants while cancelling cart", "cart_id", cartID, "err", vErr)
				return vErr
			}

			for _, item := range cart.Items {
				if item.IsCustom {
					continue
				}
				if variant, ok := variantsResp.Variants[item.VariantID]; ok && !variant.IsStockTracked() {
					continue
				}
				_, cErr = s.productProvider.UpdateVariantQuantity(gCtx,
					productparam.UpdateProductVariantQuantityAfterCancelOrderRequest{
						ProductID: item.ProductID,
						VariantID: item.VariantID,
						Quantity:  uint(item.Quantity),
					})
				if cErr != nil {
					logger.L().Error("failed to update variant quantity", "Variant_id", item.VariantID, "err", cErr)
					return cErr
				}
			}

			if cErr = s.repo.UpdateStatus(gCtx, cartID, entity.CartCanceled); err != nil {
				logger.L().Error("failed to update status", "cart_id", cartID, "err", cErr)
				return cErr
			}
			return nil
		})
	}
	_ = g.Wait()

	return cartparam.CancelExpiredPendingCartResponse{}, nil
}

// UpsertCustomItemBackoffice adds or edits a hand-typed line on an in-person
// invoice — a repair fee, an Apple ID, or a part bought from another shop —
// which exists only on this invoice and has no catalog product behind it.
func (s *BackofficeService) UpsertCustomItemBackoffice(ctx context.Context, req cartparam.UpsertCustomItemBackofficeRequest) (cartparam.UpsertCustomItemBackofficeResponse, error) {
	const op = "cartservice.UpsertCustomItemBackoffice"

	kind := req.Kind
	if kind == "" {
		kind = entity.ProductKindPhysical
	}
	// A service costs the shop nothing to hold, so carrying a unit cost on it
	// would only distort the profit report.
	unitCost := req.UnitCost
	if kind == entity.ProductKindService {
		unitCost = 0
	}

	item := entity.CartItem{
		VariantID:   req.VariantID,
		ProductName: req.Title,
		Quantity:    req.Quantity,
		PriceAtTime: req.UnitPrice,
		IsCustom:    true,
		Kind:        kind,
		UnitCost:    unitCost,
	}

	variantID, err := s.repo.UpsertCustomItem(ctx, req.UserID, item)
	if err != nil {
		return cartparam.UpsertCustomItemBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}

	return cartparam.UpsertCustomItemBackofficeResponse{VariantID: variantID}, nil
}

// SetItemManualDiscountBackoffice records the discount the salesperson gave on
// one line of an in-person invoice.
// SetItemUnitPriceBackoffice sets the unit price for a cart line. The backend
// derives manual_discount = (price_at_time - unit_price) * quantity and stores
// it, supporting both discounts (positive) and markups (negative).
func (s *BackofficeService) SetItemUnitPriceBackoffice(ctx context.Context, req cartparam.SetItemUnitPriceBackofficeRequest) (cartparam.SetItemUnitPriceBackofficeResponse, error) {
	const op = "cartservice.SetItemUnitPriceBackoffice"

	cart, err := s.repo.GetUserActiveCart(ctx, req.UserID)
	if err != nil {
		return cartparam.SetItemUnitPriceBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}

	var found *entity.CartItem
	for i := range cart.Items {
		if cart.Items[i].VariantID == req.VariantID {
			found = &cart.Items[i]
			break
		}
	}
	if found == nil {
		return cartparam.SetItemUnitPriceBackofficeResponse{}, richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("cart.item_not_found")
	}

	var adjustment int64
	if req.UnitPrice == 0 || req.UnitPrice == found.PriceAtTime {
		// Reset: clear the adjustment.
		adjustment = 0
	} else {
		adjustment = (found.PriceAtTime - req.UnitPrice) * int64(found.Quantity)
	}

	if err := s.repo.SetItemManualDiscount(ctx, req.UserID, req.VariantID, adjustment); err != nil {
		return cartparam.SetItemUnitPriceBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return cartparam.SetItemUnitPriceBackofficeResponse{}, nil
}

func (s *BackofficeService) SetItemManualDiscountBackoffice(ctx context.Context, req cartparam.SetItemManualDiscountBackofficeRequest) (cartparam.SetItemManualDiscountBackofficeResponse, error) {
	const op = "cartservice.SetItemManualDiscountBackoffice"

	if err := s.repo.SetItemManualDiscount(ctx, req.UserID, req.VariantID, int64(req.Amount)); err != nil {
		return cartparam.SetItemManualDiscountBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return cartparam.SetItemManualDiscountBackofficeResponse{}, nil
}

// SetCartManualDiscountBackoffice records the discount the salesperson gave on
// the in-person invoice as a whole.
func (s *BackofficeService) SetCartManualDiscountBackoffice(ctx context.Context, req cartparam.SetCartManualDiscountBackofficeRequest) (cartparam.SetCartManualDiscountBackofficeResponse, error) {
	const op = "cartservice.SetCartManualDiscountBackoffice"

	if err := s.repo.SetCartManualDiscount(ctx, req.UserID, req.Amount); err != nil {
		return cartparam.SetCartManualDiscountBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return cartparam.SetCartManualDiscountBackofficeResponse{}, nil
}

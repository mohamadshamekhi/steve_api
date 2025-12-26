package orderservice

import (
	"context"
	"fmt"
	"sync"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	settingparam "github.com/abolfazlnorzad/steve/param/setting"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StoreFrontService struct {
	cartProvider     StorefrontCartProvider
	discountProvider StorefrontDiscountProvider
	productProvider  StorefrontProductProvider
	repo             OrderReadWriter
	userProvider     StorefrontUserProvider
	shippingProvider ShippingProvider
}

type ShippingProvider interface {
	FindByID(ctx context.Context, req settingparam.FindByIDShippingRequest) (settingparam.FindByIDShippingResponse, error)
}

type StorefrontCartProvider interface {
	GetMyActiveCart(ctx context.Context, req cartparam.MyActiveCartRequest) (cartparam.MyActiveCartResponse, error)
	UpdateCartStatus(ctx context.Context, req cartparam.UpdateCartStatusRequest) (cartparam.UpdateCartStatusResponse, error)
}

type StorefrontDiscountProvider interface {
	GetByCode(ctx context.Context, req discountparam.GetSingleDiscountByCodeRequest) (discountparam.GetSingleDiscountByCodeResponse, error)
	GetValidDiscountsForCart(ctx context.Context, req discountparam.GetValidDiscountsForSpecificCartRequest) (discountparam.GetValidDiscountsForSpecificCartResponse, error)
	CalculateDiscountAmount(ctx context.Context, d entity.Discount, cart entity.Cart) (int64, error)
	UpdateDiscountUsageAfterRegisterOrder(ctx context.Context, req discountparam.UpdateDiscountUsageAfterRegisterOrderRequest) (discountparam.UpdateDiscountUsageAfterRegisterOrderResponse, error)
}
type StorefrontProductProvider interface {
	GetVariantsByIDs(ctx context.Context, req productparam.GetVariantsByIDsRequest) (response productparam.GetVariantsByIDsResponse, err error)
	GetProductCategoriesForCartItems(ctx context.Context, req productparam.GetProductCategoriesForCartItemsRequest) (productparam.GetProductCategoriesForCartItemsResponse, error)
	UpdateProductVariantQuantityAfterRegisterOrder(ctx context.Context, req productparam.UpdateProductVariantQuantityAfterRegisterOrderRequest) (productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse, error)
}
type StorefrontUserProvider interface {
	FindAddressByID(ctx context.Context, req userparam.FindAddressByIDRequest) (userparam.FindAddressByIDResponse, error)
}

func New(cartProvider StorefrontCartProvider,
	discountProvider StorefrontDiscountProvider,
	productProvider StorefrontProductProvider,
	repo OrderReadWriter, userProvider StorefrontUserProvider,
	shippingProvider ShippingProvider,
) StoreFrontService {
	return StoreFrontService{
		cartProvider:     cartProvider,
		discountProvider: discountProvider,
		productProvider:  productProvider,
		repo:             repo,
		userProvider:     userProvider,
		shippingProvider: shippingProvider,
	}
}

func (s StoreFrontService) RegisterOrder(ctx context.Context, req orderparam.RegisterOrderRequest) (orderparam.RegisterOrderResponse, error) {
	const op = "orderservice.GoToGateway"
	userID, _ := ctx.Value("user_id").(string)
	activeCartResp, gCErr := s.cartProvider.GetMyActiveCart(ctx, cartparam.MyActiveCartRequest{UserID: userID})
	if gCErr != nil {
		return orderparam.RegisterOrderResponse{}, richerror.New(op).WithWarpError(gCErr)
	}
	if req.DiscountCode != "" {
		discountResp, gDErr := s.discountProvider.GetByCode(ctx, discountparam.GetSingleDiscountByCodeRequest{Code: req.DiscountCode})
		if gDErr != nil {
			return orderparam.RegisterOrderResponse{}, richerror.New(op).WithWarpError(gDErr)
		}

		dErr := discountResp.Discount.CheckIsValidForUser(userID, activeCartResp.Cart.Subtotal)
		if dErr != nil {
			return orderparam.RegisterOrderResponse{}, richerror.New(op).WithWarpError(dErr)
		}
	}
	var cartItem []entity.CartItem
	for _, item := range activeCartResp.Cart.Items {
		cartItem = append(cartItem, entity.CartItem{
			ProductID:         item.ProductID,
			VariantID:         item.VariantID,
			Quantity:          item.Quantity,
			ProductName:       item.ProductName,
			VariantAttributes: item.VariantAttributes,
			PriceAtTime:       item.PriceAtTime,
		})
	}
	cart := entity.Cart{
		ID:     activeCartResp.Cart.ID,
		UserID: activeCartResp.Cart.User.ID,
		Items:  cartItem,
		Status: entity.CartStatus(activeCartResp.Cart.Status),
	}
	summary, err := s.FinalizeCartForCheckout(ctx, cart, req.DiscountCode)
	if err != nil {
		return orderparam.RegisterOrderResponse{}, richerror.New(op).WithWarpError(err)
	}

	if !summary.IsReady {
		return orderparam.RegisterOrderResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("order.cart_not_ready_for_checkout").
			WithMeta(map[string]any{"summary": summary})
	}

	add, err := s.userProvider.FindAddressByID(ctx, userparam.FindAddressByIDRequest{UserID: userID, AddressID: req.AddressID})
	if err != nil {
		return orderparam.RegisterOrderResponse{}, richerror.New(op).WithWarpError(err)
	}
	shipping, err := s.shippingProvider.FindByID(ctx, settingparam.FindByIDShippingRequest{ShippingID: req.ShippingID})
	if err != nil {
		return orderparam.RegisterOrderResponse{}, richerror.New(op).WithWarpError(err)
	}
	newOrder := entity.Order{
		UserID:           userID,
		CartID:           cart.ID,
		Items:            buildOrderItems(summary.Items),
		ShippingAddress:  add.Address,
		SubtotalAmount:   uint(summary.Subtotal),
		ShippingAmount:   0,
		AppliedDiscounts: summary.AppliedDiscounts,
		FinalAmount:      uint(summary.FinalPrice),
		Shipping:         shipping.Shipping,
		Status:           entity.OrderStatusPendingPayment,
	}
	newOrder.BeforeCreate()
	if err := s.repo.Create(ctx, newOrder); err != nil {
		return orderparam.RegisterOrderResponse{}, richerror.New(op).WithWarpError(err)
	}

	if _, err = s.productProvider.UpdateProductVariantQuantityAfterRegisterOrder(ctx, productparam.UpdateProductVariantQuantityAfterRegisterOrderRequest{Order: newOrder}); err != nil {
		return orderparam.RegisterOrderResponse{}, richerror.New(op).WithWarpError(err)
	}

	s.discountProvider.UpdateDiscountUsageAfterRegisterOrder(ctx, discountparam.UpdateDiscountUsageAfterRegisterOrderRequest{AppliedDiscount: newOrder.AppliedDiscounts})

	if _, err = s.cartProvider.UpdateCartStatus(ctx, cartparam.UpdateCartStatusRequest{
		CartID: cart.ID,
		Status: entity.PendingPayment,
	}); err != nil {
		return orderparam.RegisterOrderResponse{}, richerror.New(op).WithWarpError(err)
	}

	return orderparam.RegisterOrderResponse{OrderID: newOrder.ID}, nil
}

func buildOrderItems(checkoutItems []entity.CheckoutItemResult) []entity.OrderItem {
	orderItems := make([]entity.OrderItem, len(checkoutItems))
	for i, item := range checkoutItems {
		orderItems[i] = entity.OrderItem{
			ProductID:         item.ProductID,
			VariantID:         item.VariantID,
			Quantity:          uint(item.Quantity),
			ProductName:       item.ProductName,
			VariantAttributes: item.VariantAttributes,
			PriceAtTime:       uint(item.PriceAtTime),
		}
	}
	return orderItems
}

func (s StoreFrontService) FinalizeCartForCheckout(ctx context.Context, cart entity.Cart, discountCode string) (entity.CheckoutSummary, error) {
	const op = "order.service.FinalizeCartForCheckout"

	summary := entity.CheckoutSummary{IsReady: true}
	summary.Subtotal = int64(cart.CalculateSubtotal())

	if err := s.validateCartItems(ctx, cart, &summary); err != nil {
		return entity.CheckoutSummary{}, richerror.New(op).WithWarpError(err)
	}

	if !summary.IsReady {
		summary.FinalPrice = summary.Subtotal
		return summary, nil
	}

	validDiscountsResp, err := s.discountProvider.GetValidDiscountsForCart(ctx, discountparam.GetValidDiscountsForSpecificCartRequest{
		DiscountCode: discountCode,
		Cart:         cart,
	})
	if err != nil {
		return entity.CheckoutSummary{}, richerror.New(op).WithWarpError(err)
	}

	s.calculateFinalPrice(ctx, &cart, &summary, validDiscountsResp.Discounts, discountCode)

	return summary, nil
}

func (s StoreFrontService) validateCartItems(ctx context.Context, cart entity.Cart, summary *entity.CheckoutSummary) error {
	variantIDs := make([]string, 0, len(cart.Items))
	for _, item := range cart.Items {
		variantIDs = append(variantIDs, item.VariantID)
	}

	respProduct, err := s.productProvider.GetVariantsByIDs(ctx, productparam.GetVariantsByIDsRequest{IDs: variantIDs})
	currentVariants := respProduct.Variants
	if err != nil {
		return err
	}

	for _, item := range cart.Items {
		itemResult := entity.CheckoutItemResult{CartItem: item}
		currentVariant, ok := currentVariants[item.VariantID]
		if !ok {
			return fmt.Errorf("could not find variant %s for product %s", item.VariantID, item.ProductID)
		}

		if currentVariant.Quantity < uint64(item.Quantity) {
			itemResult.HasStockIssue = true
			summary.IsReady = false
		}
		if currentVariant.Price != uint64(item.PriceAtTime) {
			itemResult.HasPriceChanged = true
			itemResult.CurrentPrice = int64(currentVariant.Price)
			summary.IsReady = false
		} else {
			itemResult.CurrentPrice = int64(currentVariant.Price)
		}
		summary.Items = append(summary.Items, itemResult)
	}
	return nil
}

func (s StoreFrontService) applyDiscounts(cart *entity.Cart, summary *entity.CheckoutSummary, allValidDiscounts []entity.Discount, userCode string) {
	currentPrice := summary.Subtotal
	appliedDiscounts := make([]entity.AppliedDiscount, 0)
	var coupon *entity.Discount

	autoDiscounts := make([]entity.Discount, 0)
	for _, d := range allValidDiscounts {
		if d.Code != "" && d.Code == userCode {
			c := d
			coupon = &c
		} else if d.Code == "" {
			autoDiscounts = append(autoDiscounts, d)
		}
	}

	for _, discount := range autoDiscounts {
		discountAmount, err := s.discountProvider.CalculateDiscountAmount(context.Background(), discount, *cart)
		if err != nil {
			logger.L().Error("cann't apply discount - discount id : ", discount.ID)
		}
		if discountAmount > 0 {
			currentPrice -= discountAmount
			appliedDiscounts = append(appliedDiscounts, entity.AppliedDiscount{
				DiscountID:   discount.ID,
				DiscountCode: discount.Code,
				Amount:       discountAmount,
			})
		}
	}

	if coupon != nil {
		discountAmount, err := s.discountProvider.CalculateDiscountAmount(context.Background(), *coupon, *cart)
		if err != nil {
			logger.L().Error("cann't apply discount - discount id : ", coupon.ID)
		}
		if discountAmount > 0 {
			currentPrice -= discountAmount
			appliedDiscounts = append(appliedDiscounts, entity.AppliedDiscount{
				DiscountID:   coupon.ID,
				DiscountCode: coupon.Code,
				Amount:       discountAmount,
			})
		}
	}

	summary.FinalPrice = currentPrice
	summary.AppliedDiscounts = appliedDiscounts
}

func (s *StoreFrontService) calculateFinalPrice(ctx context.Context, cart *entity.Cart, summary *entity.CheckoutSummary, allValidDiscounts []entity.Discount, userCode string) error {
	const op = "orderservice.calculateFinalPrice"

	productCategories, err := s.productProvider.GetProductCategoriesForCartItems(ctx, productparam.GetProductCategoriesForCartItemsRequest{Items: cart.Items})
	if err != nil {
		return richerror.New(op).WithWarpError(err)
	}

	var autoDiscounts []entity.Discount
	var coupon *entity.Discount
	for _, d := range allValidDiscounts {
		discountCopy := d
		if d.Code != "" && d.Code == userCode {
			coupon = &discountCopy
		} else if d.Code == "" {
			autoDiscounts = append(autoDiscounts, discountCopy)
		}
	}

	var totalDiscountAmount int64 = 0
	appliedDiscountsMap := make(map[string]entity.AppliedDiscount)

	for i, item := range summary.Items {
		var itemTotalDiscount int64 = 0
		itemBasePrice := item.PriceAtTime * int64(item.Quantity)

		var bestAutoDiscountAmount int64 = 0
		var bestAutoDiscount *entity.Discount = nil
		for _, autoDiscount := range autoDiscounts {
			if s.isDiscountApplicableToItem(autoDiscount, item.CartItem, productCategories.CategoryIDs[item.ProductID], cart.UserID) {
				currentDiscountAmount := autoDiscount.CalculateAmount(uint(itemBasePrice), item.VariantID)
				if currentDiscountAmount > uint(bestAutoDiscountAmount) {
					bestAutoDiscountAmount = int64(currentDiscountAmount)
					bestAutoDiscount = &autoDiscount
				}
			}
		}

		if bestAutoDiscount != nil {
			itemTotalDiscount += bestAutoDiscountAmount
			ad := appliedDiscountsMap[bestAutoDiscount.ID]
			ad.DiscountID = bestAutoDiscount.ID
			ad.Amount += bestAutoDiscountAmount
			appliedDiscountsMap[bestAutoDiscount.ID] = ad
		}

		if coupon != nil {
			if s.isDiscountApplicableToItem(*coupon, item.CartItem, productCategories.CategoryIDs[item.ProductID], cart.UserID) {
				priceForItemCoupon := itemBasePrice - itemTotalDiscount
				couponAmountOnItem := coupon.CalculateAmount(uint(priceForItemCoupon), item.VariantID)

				itemTotalDiscount += int64(couponAmountOnItem)
				ad := appliedDiscountsMap[coupon.ID]
				ad.DiscountID = coupon.ID
				ad.DiscountCode = coupon.Code
				ad.Amount += int64(couponAmountOnItem)
				appliedDiscountsMap[coupon.ID] = ad
			}
		}

		if itemTotalDiscount > itemBasePrice {
			itemTotalDiscount = itemBasePrice
		}
		summary.Items[i].DiscountAmount = uint(itemTotalDiscount)
		totalDiscountAmount += itemTotalDiscount
	}

	summary.DiscountedSubtotal = uint(summary.Subtotal - totalDiscountAmount)
	summary.FinalPrice = int64(summary.DiscountedSubtotal)

	summary.AppliedDiscounts = make([]entity.AppliedDiscount, 0, len(appliedDiscountsMap))
	for _, ad := range appliedDiscountsMap {
		summary.AppliedDiscounts = append(summary.AppliedDiscounts, ad)
	}

	return nil
}

func (s *StoreFrontService) isDiscountApplicableToItem(
	d entity.Discount,
	item entity.CartItem,
	itemCategoryIDs []string,
	userID string,
) bool {
	// 1) فیلتر بر اساس یوزر (اگر محدود شده باشد)
	if len(d.ApplicableUserIDs) > 0 {
		userIsAllowed := false
		for _, allowedUserID := range d.ApplicableUserIDs {
			if userID == allowedUserID {
				userIsAllowed = true
				break
			}
		}
		if !userIsAllowed {
			return false
		}
	}

	// 2) تصمیم‌گیری بر اساس ApplicableTo
	switch d.ApplicableTo {

	case entity.ApplicableToProducts:
		// باید با SpecialDiscountItems چک کنیم
		if len(d.ApplicableProductIDs) == 0 {
			return false
		}

		for _, sdi := range d.ApplicableProductIDs {
			if sdi.ProductID != item.ProductID {
				continue
			}

			// اگر VariantID در تخفیف مشخص نشده باشد، یعنی برای همه ورینت‌های این محصول معتبر است
			if sdi.VariantID == nil {
				return true
			}

			// اگر VariantID مشخص شده، باید دقیقاً با Variant آیتم برابر باشد
			if *sdi.VariantID == item.VariantID {
				return true
			}
		}

		return false

	case entity.ApplicableToCategories:
		// اگر دسته‌ای تنظیم نشده، در این مود عملاً تخفیف قابل اعمال نیست
		if len(d.ApplicableCategoryIDs) == 0 {
			return false
		}

		allowedCategorySet := make(map[string]struct{}, len(d.ApplicableCategoryIDs))
		for _, catID := range d.ApplicableCategoryIDs {
			allowedCategorySet[catID] = struct{}{}
		}

		for _, catID := range itemCategoryIDs {
			if _, ok := allowedCategorySet[catID]; ok {
				return true
			}
		}

		return false

	case entity.ApplicableToUsers:
		return true

	default:
		return false
	}
}

func (s StoreFrontService) GetUserOrders(ctx context.Context, req orderparam.GetUserOrderRequest) (orderparam.GetUserOrderResponse, error) {
	const op = "orderservice.GetUserOrder"
	resp, err := s.repo.GetUserOrderWithDetails(ctx, req.UserID, req.Pagination)
	if err != nil {
		return orderparam.GetUserOrderResponse{}, richerror.New(op).WithWarpError(err)
	}

	return orderparam.GetUserOrderResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(resp.TotalCount),
		},
		Orders: resp.Orders,
	}, nil

}
func (s StoreFrontService) GetUserOrderStatics(ctx context.Context, req orderparam.GetUserOrderStaticsRequest) (orderparam.GetUserOrderStaticsResponse, error) {
	const op = "orderservice.GetUserOrderStatics"
	var pCount, cCounter, dCounter uint64

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		count, err := s.repo.GetUserOrderStatics(ctx, req.UserID, entity.OrderStatusProcessing)
		if err != nil {
			logger.L().Error("failed to get processing order statistics", "error", err, "user_id", req.UserID)
			return
		}
		pCount = count
	}()

	go func() {
		defer wg.Done()
		count, err := s.repo.GetUserOrderStatics(ctx, req.UserID, entity.OrderStatusDelivered)
		if err != nil {
			logger.L().Error("failed to get delivered order statistics", "error", err, "user_id", req.UserID)
			return
		}
		dCounter = count
	}()

	go func() {
		defer wg.Done()
		count, err := s.repo.GetUserOrderStatics(ctx, req.UserID, entity.OrderStatusRefunded)
		if err != nil {
			logger.L().Error("failed to get refunded order statistics", "error", err, "user_id", req.UserID)
			return
		}
		cCounter = count
	}()

	wg.Wait()

	return orderparam.GetUserOrderStaticsResponse{
		OrderProcessing: pCount,
		OrderDelivered:  dCounter,
		OrderRefunded:   cCounter,
	}, nil

}

package orderservice

import (
	"context"
	"fmt"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	paymentparam "github.com/abolfazlnorzad/steve/param/payment"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"golang.org/x/sync/errgroup"
)

type BackofficeService struct {
	repo             OrderReadWriter
	userProvider     UserProvider
	paymentProvider  PaymentProvider
	productProvider  ProductProvider
	discountProvider DiscountProvider
}

func NewBackofficeService(repo OrderReadWriter, userProvider UserProvider, productProvider ProductProvider, discountProvider DiscountProvider) BackofficeService {
	return BackofficeService{repo: repo, userProvider: userProvider, productProvider: productProvider, discountProvider: discountProvider}
}

func (s *BackofficeService) SetPaymentProvider(p PaymentProvider) {
	s.paymentProvider = p
}

func (s BackofficeService) GetAll(ctx context.Context, req orderparam.GetAllOrderRequest) (orderparam.GetAllOrderResponse, error) {
	const op = "orderservice.GetAll"
	resp, err := s.repo.GetAll(ctx, req.Filter, req.Pagination)
	if err != nil {
		return orderparam.GetAllOrderResponse{}, richerror.New(op).WithWarpError(err)
	}
	if resp.TotalCount == 0 {
		return orderparam.GetAllOrderResponse{
			Pagination: commonparam.PaginationResponse{
				PageSize:   req.Pagination.GetPageSize(),
				PageNumber: req.Pagination.GetPageNumber(),
				Total:      0,
			},
			Orders: nil,
		}, nil
	}
	userIDs := make([]string, 0, len(resp.Orders))
	for _, cart := range resp.Orders {
		if cart.UserID != "" {
			userIDs = append(userIDs, cart.UserID)
		}
	}

	users, err := s.userProvider.GetUsersByIDs(ctx, userparam.GetUsersByIDsRequest{IDs: userIDs})
	if err != nil {
		logger.L().Error("cann't get users , users ids: %v - and errors : %v", userIDs, err)
	}

	var list []orderparam.BackofficeOrderList

	for _, order := range resp.Orders {
		u := entity.User{}
		for _, user := range users.Users {
			if user.ID == order.UserID {
				u = entity.User{
					ID:          user.ID,
					Name:        user.Name,
					PhoneNumber: user.PhoneNumber,
					CreatedAt:   user.CreatedAt,
					UpdatedAt:   user.CreatedAt,
					Addresses:   user.Addresses,
				}
			}
		}
		payments, err := s.paymentProvider.GetPaymentsByOrderID(ctx,
			paymentparam.GetPaymentsByOrderIDRequest{OrderID: order.ID},
		)
		if err != nil {
			return orderparam.GetAllOrderResponse{}, richerror.New(op).WithWarpError(err)
		}
		list = append(list, orderparam.BackofficeOrderList{
			ID:              order.ID,
			UserPhoneNumber: u.PhoneNumber,
			UserName:        u.Name,
			UserID:          order.UserID,
			Items:           order.Items,
			ShippingAddress: order.ShippingAddress,
			SubtotalAmount:  order.SubtotalAmount,
			ShippingAmount:  order.ShippingAmount,
			FinalAmount:     order.FinalAmount,
			Payments:        payments.Payments,
			Status:          order.Status,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
			Shipping:        order.Shipping,
		})
	}

	return orderparam.GetAllOrderResponse{
		Orders: list,
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{PageSize: req.Pagination.PageSize, PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)},
	}, nil
}

func (s BackofficeService) UpdateStatus(ctx context.Context, req orderparam.UpdateStatusOrderRequest) (orderparam.UpdateStatusOrderResponse, error) {
	const op = "backofficeService.UpdateStatus"

	order, err := s.repo.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		return orderparam.UpdateStatusOrderResponse{}, nil

	}
	if !entity.CanTransitionOrderStatus(order.Status, req.Status) {
		return orderparam.UpdateStatusOrderResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("order.invalid_status_transition").
			WithMeta(map[string]any{
				"from": order.Status,
				"to":   req.Status,
			})
	}

	if err := s.repo.UpdateOrderStatus(ctx, req.OrderID, req.Status); err != nil {
		return orderparam.UpdateStatusOrderResponse{}, richerror.New(op).WithWarpError(err)
	}
	return orderparam.UpdateStatusOrderResponse{}, nil
}

func (s BackofficeService) CancelExpiredPendingOrders(ctx context.Context, req orderparam.CancelExpiredPendingOrdersRequest) (orderparam.CancelExpiredPendingOrdersResponse, error) {
	const op = "orderservice.CancelExpiredPendingOrders"

	expiredBefore := time.Now().Add(-req.CanselExpiredPendingOrders)

	orders, err := s.repo.GetExpiredPendingOrders(ctx, expiredBefore)
	if err != nil {
		return orderparam.CancelExpiredPendingOrdersResponse{}, richerror.New(op).WithWarpError(err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	for _, order := range orders {
		orderID := order.ID
		g.Go(func() error {
			var cErr error
			for _, item := range order.Items {
				_, cErr = s.productProvider.UpdateVariantQuantity(gCtx,
					productparam.UpdateProductVariantQuantityAfterCancelOrderRequest{
						ProductID: item.ProductID,
						VariantID: item.VariantID,
						Quantity:  item.Quantity,
					})
				if cErr != nil {
					logger.L().Error("failed to update variant quantity", "Variant_id", item.VariantID, "err", cErr)
					return cErr
				}
			}

			if cErr = s.repo.UpdateOrderStatus(gCtx, orderID, entity.OrderStatusCancelled); err != nil {
				logger.L().Error("failed to update status", "order_id", orderID, "err", cErr)
				return cErr
			}
			return nil
		})
	}
	_ = g.Wait()

	return orderparam.CancelExpiredPendingOrdersResponse{}, nil
}

func (s BackofficeService) FinalizeCartForCheckoutBackoffice(ctx context.Context, cart entity.Cart, discountCode string) (entity.CheckoutSummary, error) {
	const op = "order.service.FinalizeCartForCheckoutBackoffice"

	summary := entity.CheckoutSummary{IsReady: true}
	summary.Subtotal = int64(cart.CalculateSubtotal())

	if err := s.validateCartItemsBackoffice(ctx, cart, &summary); err != nil {
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

	s.calculateFinalPriceBackoffice(ctx, &cart, &summary, validDiscountsResp.Discounts, discountCode)

	return summary, nil
}

func (s BackofficeService) validateCartItemsBackoffice(ctx context.Context, cart entity.Cart, summary *entity.CheckoutSummary) error {
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

func (s *BackofficeService) calculateFinalPriceBackoffice(ctx context.Context, cart *entity.Cart, summary *entity.CheckoutSummary, allValidDiscounts []entity.Discount, userCode string) error {
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
			if s.isDiscountApplicableToItemBackoffice(autoDiscount, item.CartItem, productCategories.CategoryIDs[item.ProductID], cart.UserID) {
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
			if s.isDiscountApplicableToItemBackoffice(*coupon, item.CartItem, productCategories.CategoryIDs[item.ProductID], cart.UserID) {
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

func (s *BackofficeService) isDiscountApplicableToItemBackoffice(
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

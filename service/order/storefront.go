package orderservice

import (
	"context"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StoreFrontService struct {
	cartProvider     StorefrontCartProvider
	discountProvider StorefrontDiscountProvider
	productProvider  StorefrontProductProvider
	repo             Repository
}

type StorefrontCartProvider interface {
	GetMyActiveCart(ctx context.Context, req cartparam.MyActiveCartRequest) (cartparam.MyActiveCartResponse, error)
}

type StorefrontDiscountProvider interface {
	GetByCode(ctx context.Context, req discountparam.GetSingleDiscountByCodeRequest) (discountparam.GetSingleDiscountByCodeResponse, error)
	GetValidDiscountsForCart(ctx context.Context, req discountparam.GetValidDiscountsForSpecificCartRequest) (discountparam.GetValidDiscountsForSpecificCartResponse, error)
	CalculateDiscountAmount(ctx context.Context, d entity.Discount, cart entity.Cart) (int64, error)
}
type StorefrontProductProvider interface {
	GetVariantsByIDs(ctx context.Context, req productparam.GetVariantsByIDsRequest) (response productparam.GetVariantsByIDsResponse, err error)
}

type Repository interface {
	Create(ctx context.Context, order entity.Order) error
}

func (s StoreFrontService) RegisterOrder(ctx context.Context, req orderparam.RegisterOrderRequest) (orderparam.RegisterOrderResponse, error) {
	const op = "orderservice.GoToGateway"
	userID, _ := ctx.Value("user_id").(string)
	activeCartResp, gCErr := s.cartProvider.GetMyActiveCart(ctx, cartparam.MyActiveCartRequest{UserID: userID})
	if gCErr != nil {
		return orderparam.RegisterOrderResponse{}, richerror.New(op).WithWarpError(gCErr)
	}
	discountResp, gDErr := s.discountProvider.GetByCode(ctx, discountparam.GetSingleDiscountByCodeRequest{Code: req.DiscountCode})
	if gDErr != nil {
		return orderparam.RegisterOrderResponse{}, richerror.New(op).WithWarpError(gDErr)
	}

	dErr := discountResp.Discount.CheckIsValidForUser(userID, activeCartResp.Cart.Subtotal)
	if dErr != nil {
		return orderparam.RegisterOrderResponse{}, richerror.New(op).WithWarpError(dErr)
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

	newOrder := entity.Order{
		UserID:           userID,
		Items:            buildOrderItems(summary.Items),
		SubtotalAmount:   uint(summary.Subtotal),
		AppliedDiscounts: summary.AppliedDiscounts,
		FinalAmount:      uint(summary.FinalPrice),
		Status:           entity.OrderStatusPendingPayment,
	}
	newOrder.BeforeCreate()
	if err := s.repo.Create(ctx, newOrder); err != nil {
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

	s.applyDiscounts(&cart, &summary, validDiscountsResp.Discounts, discountCode)

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

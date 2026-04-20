package orderservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
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
	cartProvider     CartProvider
}

func NewBackofficeService(repo OrderReadWriter, userProvider UserProvider, productProvider ProductProvider, discountProvider DiscountProvider, cartProvider CartProvider) BackofficeService {
	return BackofficeService{repo: repo, userProvider: userProvider, productProvider: productProvider, discountProvider: discountProvider, cartProvider: cartProvider}
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
	const op = "orderservice.UpdateStatus"

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

func (s BackofficeService) RegisterOrderBackoffice(ctx context.Context, req orderparam.RegisterOrderBackofficeRequest) (orderparam.RegisterOrderBackofficeResponse, error) {
	const op = "orderservice.RegisterOrderBackoffice"

	myActiveCart, cErr := s.cartProvider.GetMyActiveCartBackoffice(ctx, cartparam.MyActiveCartBackofficeRequest{UserID: req.UserID})
	if cErr != nil {
		return orderparam.RegisterOrderBackofficeResponse{}, richerror.New(op).WithWarpError(cErr)
	}
	if req.DiscountCode != "" {
		dCode, pErr := s.discountProvider.GetByCode(ctx, discountparam.GetSingleDiscountByCodeRequest{Code: req.DiscountCode})
		if pErr != nil {
			return orderparam.RegisterOrderBackofficeResponse{}, richerror.New(op).WithWarpError(pErr)
		}
		if vErr := dCode.Discount.CheckIsValidForUser(myActiveCart.Cart.UserID, myActiveCart.Cart.Subtotal); vErr != nil {
			return orderparam.RegisterOrderBackofficeResponse{}, richerror.New(op).WithWarpError(vErr)
		}
	}
	var cartItem []entity.CartItem
	for _, item := range myActiveCart.Cart.Items {
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
		ID:     myActiveCart.Cart.ID,
		UserID: myActiveCart.Cart.UserID,
		Items:  cartItem,
		Status: entity.CartStatus(myActiveCart.Cart.Status),
	}
	summary, err := s.discountProvider.ApplyDiscountToCart(ctx, cart, req.DiscountCode)
	if err != nil {
		return orderparam.RegisterOrderBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	if !summary.IsReady {
		return orderparam.RegisterOrderBackofficeResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("order.cart_not_ready_for_checkout").
			WithMeta(map[string]any{"summary": summary})
	}

	newOrder := entity.Order{
		UserID:           req.UserID,
		CartID:           cart.ID,
		Items:            buildOrderItems(summary.Items),
		SubtotalAmount:   uint(summary.Subtotal),
		AppliedDiscounts: summary.AppliedDiscounts,
		FinalAmount:      uint(summary.FinalPrice),
		Status:           entity.OrderStatusPendingPayment,
	}
	newOrder.BeforeCreate()
	newOrder.PurchaseType = entity.PurchaseTypeInPerson
	if err := s.repo.Create(ctx, newOrder); err != nil {
		return orderparam.RegisterOrderBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}

	if _, err = s.cartProvider.UpdateStatus(ctx, cartparam.UpdateCartStatusRequest{CartID: cart.ID, Status: entity.PendingPayment}); err != nil {
		return orderparam.RegisterOrderBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}

	var orderItemInfo []orderparam.OrderItemInfo
	for _, item := range summary.Items {
		orderItemInfo = append(orderItemInfo, orderparam.OrderItemInfo{
			ProductID:         item.ProductID,
			VariantID:         item.VariantID,
			Quantity:          uint(item.Quantity),
			VariantAttributes: item.VariantAttributes,
			TotalPriceAtTime:  uint(item.PriceAtTime),
			DiscountAmount:    item.DiscountAmount,
		})
	}

	orderInfo := orderparam.BackofficeOrderInfo{
		ID:                  newOrder.ID,
		UserID:              newOrder.UserID,
		CartID:              newOrder.CartID,
		Items:               orderItemInfo,
		SubtotalAmount:      newOrder.SubtotalAmount,
		AppliedDiscounts:    newOrder.AppliedDiscounts,
		TotalDiscountAmount: newOrder.SubtotalAmount - summary.DiscountedSubtotal,
		FinalAmount:         newOrder.FinalAmount,
		Status:              newOrder.Status,
	}

	return orderparam.RegisterOrderBackofficeResponse{Order: orderInfo}, nil
}

func (s BackofficeService) CompleteOrderBackoffice(ctx context.Context, req orderparam.CompletedOrderBackofficeRequest) (orderparam.CompletedOrderBackofficeResponse, error) {
	const op = "orderservice.CompleteOrderBackoffice"

	order, oErr := s.repo.GetOrderByID(ctx, req.OrderID)
	if oErr != nil {
		return orderparam.CompletedOrderBackofficeResponse{}, richerror.New(op).WithWarpError(oErr)
	}

	_, err := s.paymentProvider.CreatePaymentForBackofficePanel(ctx, paymentparam.CreatePaymentForBackofficePanelRequest{OrderID: order.ID, Description: req.Description})
	if err != nil {
		return orderparam.CompletedOrderBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	if _, err = s.cartProvider.UpdateStatus(ctx, cartparam.UpdateCartStatusRequest{CartID: order.CartID, Status: entity.CartPaid}); err != nil {
		return orderparam.CompletedOrderBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	if err = s.repo.UpdateOrderStatus(ctx, order.ID, entity.OrderStatusDelivered); err != nil {
		return orderparam.CompletedOrderBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	if _, err = s.productProvider.UpdateProductVariantQuantityAfterRegisterOrder(ctx, productparam.UpdateProductVariantQuantityAfterRegisterOrderRequest{Order: order}); err != nil {
		return orderparam.CompletedOrderBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	if _, err = s.discountProvider.UpdateDiscountUsageAfterRegisterOrder(ctx, discountparam.UpdateDiscountUsageAfterRegisterOrderRequest{AppliedDiscount: order.AppliedDiscounts}); err != nil {
		return orderparam.CompletedOrderBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	if _, err = s.productProvider.UpdateProductSalesCountAfterRegisterOrder(ctx, productparam.UpdateProductSalesCountAfterRegisterOrderRequest{Order: order}); err != nil {
		return orderparam.CompletedOrderBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return orderparam.CompletedOrderBackofficeResponse{}, nil
}

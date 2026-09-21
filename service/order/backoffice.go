package orderservice

import (
	"context"
	"errors"
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
	inventory        InventoryRecorder
	cash             CashRecorder
}

func NewBackofficeService(repo OrderReadWriter, userProvider UserProvider, productProvider ProductProvider, discountProvider DiscountProvider, cartProvider CartProvider) BackofficeService {
	return BackofficeService{repo: repo, userProvider: userProvider, productProvider: productProvider, discountProvider: discountProvider, cartProvider: cartProvider}
}

func (s *BackofficeService) SetPaymentProvider(p PaymentProvider) {
	s.paymentProvider = p
}

// SetInventoryRecorder wires the stock ledger in after construction.
func (s *BackofficeService) SetInventoryRecorder(r InventoryRecorder) {
	s.inventory = r
}

// SetCashRecorder wires the cashbox in after construction.
func (s *BackofficeService) SetCashRecorder(r CashRecorder) {
	s.cash = r
}

// recordCashIn books a settled counter sale as money received. Online orders
// are settled by the payment gateway flow, not here.
func (s BackofficeService) recordCashIn(ctx context.Context, order entity.Order, method entity.PaymentMethod) {
	if s.cash == nil || order.FinalAmount == 0 {
		return
	}
	if method == "" {
		method = entity.PaymentMethodCash
	}

	t := entity.CashTransaction{
		Direction:   entity.CashIn,
		Method:      method,
		Amount:      uint64(order.FinalAmount),
		Description: "فروش حضوری",
		RefType:     entity.CashRefOrder,
		RefID:       order.ID,
		OccurredAt:  time.Now(),
	}
	if err := s.cash.Record(ctx, t); err != nil {
		logger.L().Error("failed to record order payment in cashbox", "order_id", order.ID, "err", err)
	}
}

// enrichSourcing resolves where each checkout line will be sourced from and
// what it costs. The backoffice checkout builds its summary through the
// discount service, which knows nothing about stock, so the sourcing has to
// be attached here before the order is frozen.
func (s BackofficeService) enrichSourcing(ctx context.Context, summary *entity.CheckoutSummary) {
	if len(summary.Items) == 0 {
		return
	}

	variantIDs := make([]string, 0, len(summary.Items))
	for _, item := range summary.Items {
		variantIDs = append(variantIDs, item.VariantID)
	}

	resp, err := s.productProvider.GetVariantsByIDs(ctx, productparam.GetVariantsByIDsRequest{IDs: variantIDs})
	if err != nil {
		logger.L().Error("failed to resolve variants for order sourcing", "err", err)
		return
	}

	for i := range summary.Items {
		variant, ok := resp.Variants[summary.Items[i].VariantID]
		if !ok {
			continue
		}
		applySourcing(&summary.Items[i], variant)
	}
}

// recordSale books the sale in the ledger without failing the order.
func (s BackofficeService) recordSale(ctx context.Context, order entity.Order) {
	if s.inventory == nil {
		return
	}
	if err := s.inventory.RecordOrderSale(ctx, order); err != nil {
		logger.L().Error("failed to record inventory movements for order", "order_id", order.ID, "err", err)
	}
}

// recordReturn books stock coming back from a cancelled order.
func (s BackofficeService) recordReturn(ctx context.Context, order entity.Order, item entity.OrderItem) {
	if s.inventory == nil {
		return
	}
	if err := s.inventory.RecordOrderReturn(ctx, order, item); err != nil {
		logger.L().Error("failed to record inventory return", "order_id", order.ID, "variant_id", item.VariantID, "err", err)
	}
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
			ID:                  order.ID,
			UserPhoneNumber:     u.PhoneNumber,
			UserName:            u.Name,
			UserID:              order.UserID,
			Items:               order.Items,
			ShippingAddress:     order.ShippingAddress,
			SubtotalAmount:      order.SubtotalAmount,
			ShippingAmount:      order.ShippingAmount,
			FinalAmount:         order.FinalAmount,
			AppliedDiscounts:    appliedDiscountsOrEmpty(order.AppliedDiscounts),
			TotalDiscountAmount: totalDiscountAmount(order),
			ManualDiscount:      order.ManualDiscountAmount,
			Payments:            payments.Payments,
			Status:              order.Status,
			CreatedAt:           order.CreatedAt,
			UpdatedAt:           order.UpdatedAt,
			Shipping:            order.Shipping,
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
		order := order
		orderID := order.ID
		g.Go(func() error {
			var cErr error
			for _, item := range order.Items {
				// Only what actually left the warehouse comes back: a service
				// line and an externally-sourced unit were never in stock.
				returnQty := item.StockQuantity
				if item.Supply == "" {
					returnQty = item.Quantity
				}
				if item.GetKind() == entity.ProductKindService || returnQty == 0 {
					continue
				}

				_, cErr = s.productProvider.UpdateVariantQuantity(gCtx,
					productparam.UpdateProductVariantQuantityAfterCancelOrderRequest{
						ProductID: item.ProductID,
						VariantID: item.VariantID,
						Quantity:  returnQty,
					})
				if cErr != nil {
					logger.L().Error("failed to update variant quantity", "Variant_id", item.VariantID, "err", cErr)
					return cErr
				}
				s.recordReturn(gCtx, order, item)
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
			IsCustom:          item.IsCustom,
			Kind:              item.Kind,
			UnitCost:          item.UnitCost,
			ManualDiscount:    item.ManualDiscount,
		})
	}
	// The invoice-level counter discount lives on the cart, but the panel may
	// also decide it at the moment of registration; an explicit value in the
	// request wins and is what gets recorded.
	cartManualDiscount := myActiveCart.Cart.ManualDiscount
	if req.ManualDiscountAmount != nil {
		cartManualDiscount = *req.ManualDiscountAmount
	}
	cart := entity.Cart{
		ID:             myActiveCart.Cart.ID,
		UserID:         myActiveCart.Cart.UserID,
		Items:          cartItem,
		Status:         entity.CartStatus(myActiveCart.Cart.Status),
		ManualDiscount: cartManualDiscount,
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

	s.enrichSourcing(ctx, &summary)
	applyCustomLines(&summary)
	orderManualDiscount, mErr := applyManualDiscounts(&summary, cartManualDiscount)
	if mErr != nil {
		return orderparam.RegisterOrderBackofficeResponse{}, richerror.New(op).WithWarpError(mErr)
	}

	newOrder := entity.Order{
		UserID:               req.UserID,
		CartID:               cart.ID,
		Items:                buildOrderItems(summary.Items),
		SubtotalAmount:       uint(summary.Subtotal),
		AppliedDiscounts:     summary.AppliedDiscounts,
		ManualDiscountAmount: orderManualDiscount,
		FinalAmount:          uint(summary.FinalPrice),
		Status:               entity.OrderStatusPendingPayment,
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
	var itemsManualDiscount uint
	for _, item := range summary.Items {
		itemsManualDiscount += item.ManualDiscount
		orderItemInfo = append(orderItemInfo, orderparam.OrderItemInfo{
			ProductID:         item.ProductID,
			VariantID:         item.VariantID,
			Quantity:          uint(item.Quantity),
			ProductName:       item.ProductName,
			VariantAttributes: item.VariantAttributes,
			TotalPriceAtTime:  uint(item.PriceAtTime),
			DiscountAmount:    item.DiscountAmount,
			ManualDiscount:    item.ManualDiscount,
			IsCustom:          item.IsCustom,
			Kind:              item.Kind,
		})
	}

	orderInfo := orderparam.BackofficeOrderInfo{
		ID:                  newOrder.ID,
		UserID:              newOrder.UserID,
		CartID:              newOrder.CartID,
		Items:               orderItemInfo,
		SubtotalAmount:      newOrder.SubtotalAmount,
		AppliedDiscounts:    newOrder.AppliedDiscounts,
		TotalDiscountAmount: newOrder.SubtotalAmount - newOrder.FinalAmount,
		ItemsManualDiscount: itemsManualDiscount,
		ManualDiscount:      newOrder.ManualDiscountAmount,
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
	s.recordSale(ctx, order)
	s.recordCashIn(ctx, order, req.PaymentMethod)
	if _, err = s.discountProvider.UpdateDiscountUsageAfterRegisterOrder(ctx, discountparam.UpdateDiscountUsageAfterRegisterOrderRequest{AppliedDiscount: order.AppliedDiscounts}); err != nil {
		return orderparam.CompletedOrderBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	if _, err = s.productProvider.UpdateProductSalesCountAfterRegisterOrder(ctx, productparam.UpdateProductSalesCountAfterRegisterOrderRequest{Order: order}); err != nil {
		return orderparam.CompletedOrderBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return orderparam.CompletedOrderBackofficeResponse{}, nil
}

func (s BackofficeService) CheckPendingExist(ctx context.Context, req orderparam.CheckPendingExistRequest) (orderparam.CheckPendingExistResponse, error) {
	const op = "orderservice.CheckPendingExist"

	order, err := s.repo.CheckPendingOrderExist(ctx, req.UserID)
	if err != nil {
		var rErr *richerror.RichError
		if errors.As(err, &rErr) && richerror.GetKind(rErr) == richerror.KindNotFound {
			return orderparam.CheckPendingExistResponse{}, nil
		}
		return orderparam.CheckPendingExistResponse{}, richerror.New(op).WithWarpError(err)
	}

	var itemsInfo []orderparam.OrderItemInfo
	for _, item := range order.Items {
		itemsInfo = append(itemsInfo, orderparam.OrderItemInfo{
			ProductID:         item.ProductID,
			VariantID:         item.VariantID,
			Quantity:          item.Quantity,
			VariantAttributes: item.VariantAttributes,
			TotalPriceAtTime:  item.PriceAtTime,
		})
	}

	orderInfo := orderparam.BackofficeOrderInfo{
		ID:               order.ID,
		UserID:           order.UserID,
		CartID:           order.CartID,
		Items:            itemsInfo,
		AppliedDiscounts: order.AppliedDiscounts,
		SubtotalAmount:   order.SubtotalAmount,
		FinalAmount:      order.FinalAmount,
		Status:           order.Status,
	}

	return orderparam.CheckPendingExistResponse{
		Order: &orderInfo,
	}, err
}

// applyCustomLines freezes sourcing onto the hand-typed lines of an in-person
// invoice. They never passed through the catalog, so enrichSourcing has
// nothing to resolve them from: a service line ships nothing, and a part
// bought from another shop is external supply costing whatever the
// salesperson typed.
func applyCustomLines(summary *entity.CheckoutSummary) {
	for i := range summary.Items {
		item := &summary.Items[i]
		if !item.IsCustom {
			continue
		}
		item.Kind = item.CartItem.GetKind()
		item.CostAtTime = item.UnitCost
		if item.Kind == entity.ProductKindService {
			item.Supply = entity.SupplySourceService
			continue
		}
		item.Supply = entity.SupplySourceExternal
		item.ExternalQuantity = uint(item.Quantity)
	}
}

// applyManualDiscounts subtracts what was given away by hand at the counter:
// first the per-line discounts, on top of whatever campaign discount each line
// already earned, then one discount on the invoice as a whole.
//
// A discount larger than what is left to pay is refused rather than clamped:
// it means the salesperson mistyped, or a campaign already took the line down,
// and silently charging a different amount than the panel showed would be
// worse than making them re-enter it.
func applyManualDiscounts(summary *entity.CheckoutSummary, invoiceLevel uint) (uint, error) {
	const op = "orderservice.applyManualDiscounts"

	var totalLineDiscount int64
	for i := range summary.Items {
		item := &summary.Items[i]

		manual := int64(item.ManualDiscount)
		if manual <= 0 {
			item.ManualDiscount = 0
			continue
		}

		remaining := item.LineTotal() - int64(item.DiscountAmount)
		if manual > remaining {
			return 0, richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("order.manual_discount_exceeds_total").
				WithMeta(map[string]any{
					"variant_id": item.VariantID,
					"line_total": item.LineTotal(),
					"payable":    remaining,
					"requested":  manual,
				})
		}

		item.DiscountAmount += uint(manual)
		totalLineDiscount += manual
	}

	summary.FinalPrice -= totalLineDiscount

	applied := int64(invoiceLevel)
	if applied > summary.FinalPrice {
		return 0, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("order.manual_discount_exceeds_total").
			WithMeta(map[string]any{
				"payable":   summary.FinalPrice,
				"requested": applied,
			})
	}
	summary.FinalPrice -= applied

	return uint(applied), nil
}

// totalDiscountAmount is what the invoice lost to discounts of any kind.
// FinalAmount already includes shipping, so it is added back before diffing.
func totalDiscountAmount(o entity.Order) uint {
	gross := o.SubtotalAmount + o.ShippingAmount
	if o.FinalAmount >= gross {
		return 0
	}
	return gross - o.FinalAmount
}

// appliedDiscountsOrEmpty keeps the field a JSON array for the frontend
// instead of null on orders with no discount.
func appliedDiscountsOrEmpty(d []entity.AppliedDiscount) []entity.AppliedDiscount {
	if d == nil {
		return []entity.AppliedDiscount{}
	}
	return d
}

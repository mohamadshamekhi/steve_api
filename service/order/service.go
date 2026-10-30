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
)

type ListResult struct {
	Orders     []entity.Order
	TotalCount int64
}

type StorefrontListResult struct {
	Orders     []orderparam.StoreFrontOrderInfo
	TotalCount int64
}

type OrderReader interface {
	GetAll(ctx context.Context, f orderparam.BackofficeFilterOrder, p commonparam.PaginationRequest) (ListResult, error)
	GetUserOrderWithDetails(ctx context.Context, UserID string, p commonparam.PaginationRequest) (StorefrontListResult, error)
	GetUserOrderStatics(ctx context.Context, userID string, status entity.OrderStatus) (uint64, error)
	GetOrderByID(ctx context.Context, orderID string) (entity.Order, error)
	GetExpiredPendingOrders(ctx context.Context, expiredBefore time.Time) ([]entity.Order, error)
	CheckPendingOrderExist(ctx context.Context, userID string) (*entity.Order, error)
}

type OrderWriter interface {
	Create(ctx context.Context, order entity.Order) error
	UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error
	SetCancelReason(ctx context.Context, orderID string, reason entity.CancelReason) error
	SetTrackingCode(ctx context.Context, orderID string, trackingCode string) error
}

type OrderReadWriter interface {
	OrderReader
	OrderWriter
}

type UserProvider interface {
	GetUsersByIDs(ctx context.Context, req userparam.GetUsersByIDsRequest) (userparam.GetUsersByIDsResponse, error)
}

type PaymentProvider interface {
	GetPaymentsByOrderID(ctx context.Context, req paymentparam.GetPaymentsByOrderIDRequest) (paymentparam.GetPaymentsByOrderIDResponse, error)
	CreatePaymentForBackofficePanel(ctx context.Context, req paymentparam.CreatePaymentForBackofficePanelRequest) (paymentparam.CreatePaymentForBackofficePanelResponse, error)
}

type CartProvider interface {
	GetMyActiveCartBackoffice(ctx context.Context, req cartparam.MyActiveCartBackofficeRequest) (cartparam.MyActiveCartBackofficeResponse, error)
	UpdateStatus(ctx context.Context, req cartparam.UpdateCartStatusRequest) (cartparam.UpdateCartStatusResponse, error)
}

type ProductProvider interface {
	UpdateVariantQuantity(ctx context.Context, req productparam.UpdateProductVariantQuantityAfterCancelOrderRequest) (productparam.UpdateProductVariantQuantityAfterCancelOrderResponse, error)
	GetVariantsByIDs(ctx context.Context, req productparam.GetVariantsByIDsRequest) (response productparam.GetVariantsByIDsResponse, err error)
	GetProductCategoriesForCartItems(ctx context.Context, req productparam.GetProductCategoriesForCartItemsRequest) (productparam.GetProductCategoriesForCartItemsResponse, error)
	UpdateProductVariantQuantityAfterRegisterOrder(ctx context.Context, req productparam.UpdateProductVariantQuantityAfterRegisterOrderRequest) (productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse, error)
	UpdateProductSalesCountAfterRegisterOrder(ctx context.Context, req productparam.UpdateProductSalesCountAfterRegisterOrderRequest) (productparam.UpdateProductSalesCountAfterRegisterOrderResponse, error)
}

// InventoryRecorder writes the accounting/inventory ledger for a sale. It is
// deliberately fire-and-forget at the call sites: an order is a fact even if
// its bookkeeping line fails to land, and the ledger can be rebuilt from the
// orders themselves.
type InventoryRecorder interface {
	RecordOrderSale(ctx context.Context, order entity.Order) error
	RecordOrderReturn(ctx context.Context, order entity.Order, item entity.OrderItem) error
}

// CashRecorder books money hitting the till when a counter sale is settled.
// Like InventoryRecorder it is optional and fire-and-forget: the sale stands
// even if its cash line does not land.
type CashRecorder interface {
	Record(ctx context.Context, t entity.CashTransaction) error
}

type DiscountProvider interface {
	GetByCode(ctx context.Context, req discountparam.GetSingleDiscountByCodeRequest) (discountparam.GetSingleDiscountByCodeResponse, error)
	GetValidDiscountsForCart(ctx context.Context, req discountparam.GetValidDiscountsForSpecificCartRequest) (discountparam.GetValidDiscountsForSpecificCartResponse, error)
	CalculateDiscountAmount(ctx context.Context, d entity.Discount, cart entity.Cart) (int64, error)
	ApplyDiscountToCart(ctx context.Context, cart entity.Cart, discountCode string) (entity.CheckoutSummary, error)
	UpdateDiscountUsageAfterRegisterOrder(ctx context.Context, req discountparam.UpdateDiscountUsageAfterRegisterOrderRequest) (discountparam.UpdateDiscountUsageAfterRegisterOrderResponse, error)
}

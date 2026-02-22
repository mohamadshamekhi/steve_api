package orderservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
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
}

type OrderWriter interface {
	Create(ctx context.Context, order entity.Order) error
	UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error
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
}

type ProductProvider interface {
	UpdateVariantQuantity(ctx context.Context, req productparam.UpdateProductVariantQuantityAfterCancelOrderRequest) (productparam.UpdateProductVariantQuantityAfterCancelOrderResponse, error)
	GetVariantsByIDs(ctx context.Context, req productparam.GetVariantsByIDsRequest) (response productparam.GetVariantsByIDsResponse, err error)
	GetProductCategoriesForCartItems(ctx context.Context, req productparam.GetProductCategoriesForCartItemsRequest) (productparam.GetProductCategoriesForCartItemsResponse, error)
	UpdateProductVariantQuantityAfterRegisterOrder(ctx context.Context, req productparam.UpdateProductVariantQuantityAfterRegisterOrderRequest) (productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse, error)
}

type DiscountProvider interface {
	GetByCode(ctx context.Context, req discountparam.GetSingleDiscountByCodeRequest) (discountparam.GetSingleDiscountByCodeResponse, error)
	GetValidDiscountsForCart(ctx context.Context, req discountparam.GetValidDiscountsForSpecificCartRequest) (discountparam.GetValidDiscountsForSpecificCartResponse, error)
	CalculateDiscountAmount(ctx context.Context, d entity.Discount, cart entity.Cart) (int64, error)
	UpdateDiscountUsageAfterRegisterOrder(ctx context.Context, req discountparam.UpdateDiscountUsageAfterRegisterOrderRequest) (discountparam.UpdateDiscountUsageAfterRegisterOrderResponse, error)
}

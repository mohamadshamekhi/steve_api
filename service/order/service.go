package orderservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
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

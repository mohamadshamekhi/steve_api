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

type OrderReader interface {
	GetAll(ctx context.Context, f orderparam.BackofficeFilterOrder, p commonparam.PaginationRequest) (ListResult, error)
}

type OrderWriter interface {
	Create(ctx context.Context, order entity.Order) error
}

type OrderReadWriter interface {
	OrderReader
	OrderWriter
}

type UserProvider interface {
	GetUsersByIDs(ctx context.Context, req userparam.GetUsersByIDsRequest) (userparam.GetUsersByIDsResponse, error)
}

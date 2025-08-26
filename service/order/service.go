package orderservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
)

type OrderReader interface {
	FindByID(ctx context.Context, id string) (entity.Order, error)
}

type OrderWriter interface {
	AddOrder(ctx context.Context, cart entity.Cart) error
}

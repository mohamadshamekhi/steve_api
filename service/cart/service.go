package cartservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
)

type CartReader interface {
	FindByID(ctx context.Context, id string) (entity.Cart, error)
}

type CartWriter interface {
	AddToCart(ctx context.Context, userID string, item entity.CartItem) error
	RemoveFromCart(ctx context.Context, userID string, productID string, variantID string) error
	//Update(ctx context.Context, cart entity.Cart) error
}

type CartReadWriter interface {
	CartWriter
	CartReader
}

type PRepository interface {
	FindByID(ctx context.Context, id string) (entity.Product, error)
}

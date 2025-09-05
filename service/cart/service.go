package cartservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type CartReader interface {
	FindByID(ctx context.Context, id string) (entity.Cart, error)
	GetUserActiveCart(ctx context.Context, userID string) (*entity.Cart, error)
}

type CartWriter interface {
	AddToCart(ctx context.Context, userID string, item entity.CartItem) error
	RemoveFromCart(ctx context.Context, userID string, productID string, variantID string) error
	UpdateStatus(ctx context.Context, cartID string, status entity.CartStatus) error
}
type ListResult struct {
	Carts      []entity.Cart
	TotalCount int64
}

type CartBackofficeWrite interface {
	CartWriter
	GetAll(ctx context.Context, filter cartparam.BackofficeFilterCart, pagination commonparam.PaginationRequest) (ListResult, error)
}

type CartReadWriter interface {
	CartWriter
	CartReader
}
type CartBackofficeReadWriter interface {
	CartBackofficeWrite
	CartReadWriter
}

type PRepository interface {
	FindByID(ctx context.Context, id string) (entity.Product, error)
}

type URepository interface {
	FindManyByIDs(ctx context.Context, ids []string) ([]entity.User, error)
}

package cartservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	productparam "github.com/abolfazlnorzad/steve/param/product"
)

type CartReader interface {
	FindByID(ctx context.Context, id string) (entity.Cart, error)
	GetUserActiveCart(ctx context.Context, userID string) (*entity.Cart, error)
	GetExpiredPendingCart(ctx context.Context, expiredBefore time.Time) ([]entity.Cart, error)
}

type CartWriter interface {
	AddToCart(ctx context.Context, userID string, item entity.CartItem) error
	RemoveFromCart(ctx context.Context, userID string, productID string, variantID string) error
	UpdateStatus(ctx context.Context, cartID string, status entity.CartStatus) error
	UpdateItemPrice(ctx context.Context, userID string, variantID string, newPrice int64) error
	DeleteItemFromCart(ctx context.Context, userID string, variantID string) error
	ClearCartItemsByCartID(ctx context.Context, cartID string) error
	CreateCart(ctx context.Context, cart entity.Cart) error
}
type ListResult struct {
	Carts      []entity.Cart
	TotalCount int64
}

type CartBackofficeWrite interface {
	CartWriter
	// UpsertCustomItem adds or rewrites a hand-typed in-person invoice line
	// and returns the id that identifies it inside the cart.
	UpsertCustomItem(ctx context.Context, userID string, item entity.CartItem) (string, error)
	SetItemManualDiscount(ctx context.Context, userID string, variantID string, amount int64) error
	SetCartManualDiscount(ctx context.Context, userID string, amount uint) error
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
	GetVariantsByIDs(ctx context.Context, variantIDs []string) (map[string]entity.Variant, error)
	Delete(ctx context.Context, product entity.Product) error
}

type URepository interface {
	FindManyByIDs(ctx context.Context, ids []string) ([]entity.User, error)
	GetUserByID(ctx context.Context, id string) (entity.User, error)
}

type ProductProvider interface {
	UpdateVariantQuantity(ctx context.Context, req productparam.UpdateProductVariantQuantityAfterCancelOrderRequest) (productparam.UpdateProductVariantQuantityAfterCancelOrderResponse, error)
	GetVariantsByIDs(ctx context.Context, req productparam.GetVariantsByIDsRequest) (productparam.GetVariantsByIDsResponse, error)
}

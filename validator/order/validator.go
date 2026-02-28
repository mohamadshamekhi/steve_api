package ordervalidator

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
)

type Validator struct {
	discountProvider DiscountProvider
	cartRepo         CartRepo
}

type DiscountProvider interface {
	GetByCode(ctx context.Context, req discountparam.GetSingleDiscountByCodeRequest) (discountparam.GetSingleDiscountByCodeResponse, error)
	GetValidDiscountsForCart(ctx context.Context, req discountparam.GetValidDiscountsForSpecificCartRequest) (discountparam.GetValidDiscountsForSpecificCartResponse, error)
}

type CartRepo interface {
	GetUserActiveCart(ctx context.Context, userID string) (*entity.Cart, error)
}

func New(discountProvider DiscountProvider, cartRepo CartRepo) Validator {
	return Validator{discountProvider: discountProvider, cartRepo: cartRepo}
}

package ordervalidator

import (
	"context"

	discountparam "github.com/abolfazlnorzad/steve/param/discount"
)

type Validator struct {
	discountProvider DiscountProvider
}

type DiscountProvider interface {
	GetByCode(ctx context.Context, req discountparam.GetSingleDiscountByCodeRequest) (discountparam.GetSingleDiscountByCodeResponse, error)
}

func New(discountProvider DiscountProvider) Validator {
	return Validator{discountProvider: discountProvider}
}

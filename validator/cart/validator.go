package cartvalidator

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Validator struct {
	productRepo PRepository
	repo        Repository
}

type Repository interface {
	FindByID(ctx context.Context, id string) (entity.Cart, error)
}

type PRepository interface {
	FindByID(ctx context.Context, id string) (entity.Product, error)
}

func New(pr PRepository, repo Repository) Validator {
	return Validator{
		productRepo: pr,
		repo:        repo,
	}
}

func (v Validator) checkProductIsValid(value any) error {
	const op = "cartvalidator.checkProductIsValid"
	var id string
	switch v := value.(type) {
	case string:
		id = v
	case *string:
		if v == nil {
			return nil
		}
		id = *v
	default:
		return richerror.New("brandvalidator.BrandIsUnique.type_assertion").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type for brand name"))
	}
	p, err := v.productRepo.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	if !p.IsPublished {
		return richerror.New(op).WithKind(richerror.KindBadRequest).WithMessageKey("product.is-not-published")
	}
	return nil
}
func (v Validator) checkVariantIsValid(productID string) validation.RuleFunc {
	return func(value any) error {
		const op = "cartvalidator.checkVariantIsValid"
		p, _ := v.productRepo.FindByID(context.Background(), productID)
		for _, variant := range p.Variants {
			if variant.ID == value {
				if variant.Quantity > 0 {
					return nil
				}
				return richerror.New(op).WithKind(richerror.KindBadRequest).WithMessageKey("product.sold-out")

			}
		}
		return richerror.New(op).WithKind(richerror.KindBadRequest).WithMessageKey("global.not-found")
	}
}

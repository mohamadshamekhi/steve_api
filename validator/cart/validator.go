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
	userRepo    URepository
}

type Repository interface {
	FindByID(ctx context.Context, id string) (entity.Cart, error)
}

type PRepository interface {
	FindByID(ctx context.Context, id string) (entity.Product, error)
}

type URepository interface {
	GetUserByID(ctx context.Context, id string) (entity.User, error)
}

func New(pr PRepository, repo Repository, uRepo URepository) Validator {
	return Validator{
		productRepo: pr,
		repo:        repo,
		userRepo:    uRepo,
	}
}

// productIDFromValue unwraps the string / *string shapes ozzo hands a rule.
func productIDFromValue(value any) (string, bool, error) {
	switch v := value.(type) {
	case string:
		return v, true, nil
	case *string:
		if v == nil {
			return "", false, nil
		}
		return *v, true, nil
	default:
		return "", false, richerror.New("cartvalidator.productIDFromValue").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type for product id"))
	}
}

// checkProductIsSellableOnWebsite gates the storefront cart: the product must
// be published, not deleted and opted into the website channel. An in-shop-only
// item must never reach a customer's cart even if its ID is guessed.
func (v Validator) checkProductIsSellableOnWebsite(value any) error {
	const op = "cartvalidator.checkProductIsSellableOnWebsite"

	id, ok, err := productIDFromValue(value)
	if err != nil || !ok {
		return err
	}

	p, err := v.productRepo.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	if !p.IsSellableOnWebsite() {
		return richerror.New(op).WithKind(richerror.KindBadRequest).
			WithMessageKey("product.is-not-published")
	}
	return nil
}

// checkProductIsSellableInPerson gates the backoffice cart, which writes
// over-the-counter invoices. It deliberately ignores IsPublished: a service or
// a shop-only item has no storefront page by design, yet must be invoiceable.
func (v Validator) checkProductIsSellableInPerson(value any) error {
	const op = "cartvalidator.checkProductIsSellableInPerson"

	id, ok, err := productIDFromValue(value)
	if err != nil || !ok {
		return err
	}

	p, err := v.productRepo.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	if !p.IsSellableInPerson() {
		return richerror.New(op).WithKind(richerror.KindBadRequest).
			WithMessageKey("product.not_sellable_in_person")
	}
	return nil
}

// checkVariantIsValid reports whether one unit of the variant can be sold.
//
// It asks the variant itself rather than testing Quantity > 0: an on-demand
// variant is sourced per order and holds no stock, a backorderable one keeps
// selling past zero, and a service has no stock at all — all three are
// sellable while showing a quantity of zero.
func (v Validator) checkVariantIsValid(productID string) validation.RuleFunc {
	return func(value any) error {
		const op = "cartvalidator.checkVariantIsValid"

		p, _ := v.productRepo.FindByID(context.Background(), productID)
		for _, variant := range p.Variants {
			if variant.ID == value {
				if variant.CanSell(1) {
					return nil
				}
				return richerror.New(op).WithKind(richerror.KindBadRequest).
					WithMessageKey("product.sold-out")
			}
		}
		return richerror.New(op).WithKind(richerror.KindBadRequest).
			WithMessageKey("global.not-found")
	}
}

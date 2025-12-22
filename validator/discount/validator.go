package discountvalidator

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type Validator struct {
	repo             Repository
	categoryProvider CategoryProvider
	productProvider  ProductProvider
	userProvider     UserProvider
}

func New(repo Repository, categoryProvider CategoryProvider, productProvider ProductProvider, userProvider UserProvider) Validator {
	return Validator{
		repo:             repo,
		categoryProvider: categoryProvider,
		productProvider:  productProvider,
		userProvider:     userProvider,
	}
}

type CategoryProvider interface {
	FindNonExistingIDs(ctx context.Context, ids []string) ([]string, error)
}
type ProductProvider interface {
	FindNonExistingIDs(ctx context.Context, ids []string) ([]string, error)
}
type UserProvider interface {
	FindNonExistingIDs(ctx context.Context, ids []string) ([]string, error)
}

type Repository interface {
	CheckCodeIsUnique(ctx context.Context, code string) (bool, error)
	FindByID(ctx context.Context, id string) (entity.Discount, error)
}

func (v Validator) checkCodeIsUnique(value any) error {
	code, _ := value.(string)
	resp, err := v.repo.CheckCodeIsUnique(context.Background(), code)
	if err != nil {
		return err
	}
	if resp {
		return nil
	}
	return errors.New("code is unique")
}

func (v Validator) checkDiscountIsExist(value any) error {
	const op = "discountvalidator.checkDiscountIsExist"
	id, ok := value.(string)
	if !ok {
		return richerror.New(op).WithKind(richerror.KindBadRequest).WithMessageKey("global.validation_failed")
	}
	d, err := v.repo.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	if d.ID == "" {
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
	}

	return nil
}

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

func (v Validator) categoriesExist(value any) error {
	const op = "mongodbcategory.categoriesExist"

	var list []string
	switch v := value.(type) {
	case []string:
		list = v
	case *[]string:
		if v == nil {
			return nil
		}
		list = *v
	default:
		return richerror.New(op).
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type for brand name"))
	}

	notFoundIDs, err := v.categoryProvider.FindNonExistingIDs(context.Background(), list)
	if err != nil {
		return err
	}
	if len(notFoundIDs) == 0 {
		return nil
	}
	return richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("validation.invalid_ids_exist").WithMeta(map[string]any{"InvalidIDs": notFoundIDs})
}

func (v Validator) productsExist(value any) error {
	const op = "discountvalidator.productsExist"
	var list []string
	switch v := value.(type) {
	case []string:
		list = v
	case *[]string:
		if v == nil {
			return nil
		}
		list = *v
	default:
		return richerror.New(op).
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type for brand name"))
	}
	notFoundIDs, err := v.productProvider.FindNonExistingIDs(context.Background(), list)
	if err != nil {
		return err
	}
	if len(notFoundIDs) == 0 {
		return nil
	}
	return richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("validation.invalid_ids_exist").WithMeta(map[string]any{"InvalidIDs": notFoundIDs})
}

func (v Validator) usersExist(value any) error {
	const op = "discountvalidator.usersExist"
	var list []string
	switch v := value.(type) {
	case []string:
		list = v
	case *[]string:
		if v == nil {
			return nil
		}
		list = *v
	default:
		return richerror.New(op).
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type for brand name"))
	}
	notFoundIDs, err := v.userProvider.FindNonExistingIDs(context.Background(), list)
	if err != nil {
		return err
	}
	if len(notFoundIDs) == 0 {
		return nil
	}
	return richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("validation.invalid_ids_exist").WithMeta(map[string]any{"InvalidIDs": notFoundIDs})
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

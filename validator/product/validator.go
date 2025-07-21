package productvalidator

import (
	"context"
	"errors"
	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
)

type Validator struct {
	categoryProvider CategoryProvider
	brandProvider    BrandProvider
	repo             Repository
	mediaSvc         mediaservice.Service
}

type CategoryProvider interface {
	GetCategory(ctx context.Context, req categoryparam.GetCategoryRequest) (categoryparam.GetCategoryResponse, error)
}

type BrandProvider interface {
	GetBrand(ctx context.Context, req brandparam.GetBrandRequest) (brandparam.GetBrandResponse, error)
}

type Repository interface {
	FindByID(ctx context.Context, id string) (entity.Product, error)
}

func New(CategoryProvider CategoryProvider, BrandProvider BrandProvider, r Repository, mediaSvc mediaservice.Service) Validator {
	return Validator{
		categoryProvider: CategoryProvider,
		brandProvider:    BrandProvider,
		repo:             r,
		mediaSvc:         mediaSvc,
	}
}

func (v Validator) checkBrandIsExist(value any) error {
	id, ok := value.(string)
	if !ok {
		return richerror.New("productvalidator.checkBrandIsExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	_, err := v.brandProvider.GetBrand(context.Background(), brandparam.GetBrandRequest{ID: id})
	if err != nil {
		var rErr *richerror.RichError
		if errors.As(err, &rErr) {
			if rErr.Kind() == richerror.KindNotFound {
				return errors.New(rErr.MessageKey())
			}
		}
		return err
	}

	return nil
}

func (v Validator) checkCategoryIsExist(value any) error {
	id, ok := value.(string)
	if !ok {
		return richerror.New("productvalidator.checkCategoryIsExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	resp, err := v.categoryProvider.GetCategory(context.Background(), categoryparam.GetCategoryRequest{ID: id})
	if err != nil {
		var rErr *richerror.RichError
		if errors.As(err, &rErr) {
			if rErr.Kind() == richerror.KindNotFound {
				return errors.New(rErr.MessageKey())
			}
		}
		return err
	}
	if resp.Category.ID == "" {
		return errors.New("global.not-found")
	}
	return nil
}

func (v Validator) checkProductIsExist(value any) error {
	id, ok := value.(string)
	if !ok {
		return richerror.New("productvalidator.checkProductIsExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	_, err := v.repo.FindByID(context.Background(), id)
	if err != nil {
		var rErr *richerror.RichError
		if errors.As(err, &rErr) {
			if rErr.Kind() == richerror.KindNotFound {
				return errors.New(rErr.MessageKey())
			}
		}
		return err
	}
	return nil
}

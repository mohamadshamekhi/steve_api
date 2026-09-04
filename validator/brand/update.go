package brandvalidator

import (
	"context"
	"errors"

	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateBrandRequest(req brandparam.UpdateBrandRequest) error {
	const op = "validator.ValidateUpdateBrandRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.BrandID,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Name,
			validation.When(req.Name != nil,
				validation.Required.Error("validation.required"),
				validation.By(v.brandNameIsUniqueForUpdate(req.BrandID)),
			)),
		validation.Field(&req.Slug,
			validation.When(req.Slug != nil && *req.Slug != "",
				validation.By(v.brandSlugIsUniqueForUpdate(req.BrandID)),
			),
		),
		validation.Field(&req.Description,
			validation.When(req.Description != nil,
				validation.Required.Error("validation.required"),
			),
		),
		validation.Field(&req.MediaID,
			validation.When(req.MediaID != nil,
				validation.Required.Error("validation.required"),
				validation.By(v.mediaBrandIsExist),
			),
		),
		validation.Field(&req.IsFavorite,
			validation.When(req.IsFavorite != nil,
				validation.Required.Error("validation.required"),
			)),
	); err != nil {
		var vErrs validation.Errors
		if errors.As(err, &vErrs) {
			return richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("global.validation_failed").
				WithMeta(map[string]any{"fields": vErrs})
		}

		return richerror.New(op).WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong").WithWarpError(err)
	}

	//nolint
	return nil
}

func (v Validator) brandSlugIsUniqueForUpdate(currentBrandID string) validation.RuleFunc {
	return func(value interface{}) error {
		slug, ok := value.(*string)
		if !ok || slug == nil || *slug == "" {
			return nil
		}

		existingBrand, err := v.repo.FindBySlug(context.Background(), *slug)
		if err != nil {
			var rErr *richerror.RichError
			if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
				return nil
			}
			return err
		}

		if existingBrand.ID != currentBrandID {
			return errors.New("validation.brand_slug_is_already_exist")
		}

		return nil
	}
}

func (v Validator) brandNameIsUniqueForUpdate(currentBrandID string) validation.RuleFunc {
	return func(value interface{}) error {
		name, ok := value.(*string)
		if !ok || name == nil || *name == "" {
			return nil
		}

		existingBrand, err := v.repo.FindByName(context.Background(), *name)
		if err != nil {
			var rErr *richerror.RichError
			if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
				return nil
			}
			return err
		}

		if existingBrand.ID != currentBrandID {
			return errors.New("validation.brand_is_already_exist")
		}

		return nil
	}
}

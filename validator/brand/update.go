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
		validation.Field(&req.Description,
			validation.When(req.Name != nil,
				validation.Required.Error("validation.required"),
			),
		),
		validation.Field(&req.MediaID,
			validation.When(req.MediaID != nil,
				validation.Required.Error("validation.required"),
				validation.By(v.mediaBrandIsExist),
			),
		),
	); err != nil {
		var vErrs validation.Errors
		if errors.As(err, &vErrs) {
			fieldErrorKeys := make(map[string]string)
			for field, errMessage := range vErrs {
				if errMessage != nil {
					fieldErrorKeys[field] = errMessage.Error()
				}
			}

			return richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("global.validation_failed").
				WithMeta(map[string]any{"fields": fieldErrorKeys})
		}

		return richerror.New(op).WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong").WithWarpError(err)
	}

	//nolint
	return nil
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

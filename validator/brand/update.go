package brandvalidator

import (
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
				validation.By(v.brandIsExist),
				validation.By(v.brandIsUnique),
			),
		),
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

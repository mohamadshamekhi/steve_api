package productvalidator

import (
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateSpecificationsRequest(req productparam.AddProductAttributeRequest) error {
	const op = "productvalidator.ValidateUpdateSpecificationsRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Attrs, validation.Required.Error("validation.required")),
		validation.Field(&req.ProductID, validation.Required.Error("validation.required"), validation.By(v.checkProductIsExist)),

		validation.Field(&req.Attrs, validation.Each(
			// هر آیتم باید یک struct معتبر باشد.
			validation.By(func(value interface{}) error {
				spec, ok := value.(entity.Specification)
				if !ok {
					return errors.New("invalid specification item")
				}
				return validation.ValidateStruct(&spec,
					validation.Field(&spec.Key, validation.Required.Error("validation.spec_key_required")),
					validation.Field(&spec.Value, validation.Required.Error("validation.spec_value_required")),
				)
			}),
		)),
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
		return richerror.New(op).WithKind(richerror.KindUnexpected).WithWarpError(err)
	}

	return nil
}

func (v Validator) ValidateUpdateSingleSpecificationRequest(req productparam.UpdateSingleAttributeRequest) error {
	const op = "productvalidator.ValidateUpdateSingleSpecificationRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Key, validation.Required.Error("validation.required")),
		validation.Field(&req.Value, validation.Required.Error("validation.required")),
		validation.Field(&req.ProductID, validation.Required.Error("validation.required"), validation.By(v.checkProductIsExist)),
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
		return richerror.New(op).WithKind(richerror.KindUnexpected).WithWarpError(err)
	}

	return nil
}

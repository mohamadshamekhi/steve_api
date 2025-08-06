package productvalidator

import (
	"errors"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateAddVariantRequest(req productparam.AddVariantRequest) error {
	const op = "productvalidator.ValidateAddVariantRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Price, validation.Required.Error("validation.required")),
		validation.Field(&req.Key, validation.Required.Error("validation.required")),
		validation.Field(&req.Value, validation.Required.Error("validation.required")),
		validation.Field(&req.Label, validation.Required.Error("validation.required")),
		validation.Field(&req.Warranty, validation.Required.Error("validation.required")),
		validation.Field(&req.Quantity, validation.Required.Error("validation.required")),
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

func (v Validator) ValidateUpdateVariantRequest(req productparam.UpdateVariantRequest) error {
	const op = "productvalidator.ValidateUpdateVariantRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Price, validation.When(req.Price != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Key, validation.When(req.Key != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Value, validation.When(req.Value != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Label, validation.When(req.Label != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Warranty, validation.When(req.Warranty != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Quantity, validation.When(req.Quantity != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.ProductID, validation.Required.Error("validation.required"), validation.By(v.checkProductIsExist)),
		validation.Field(&req.VariantID, validation.Required.Error("validation.required")),
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

func (v Validator) ValidateDeleteVariantRequest(req productparam.DeleteVariantRequest) error {
	const op = "productvalidator.ValidateDeleteVariantRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.VariantID, validation.Required.Error("validation.required")),
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

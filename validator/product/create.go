package productvalidator

import (
	"errors"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateCreateProductRequest(req productparam.CreateProductRequest) error {
	const op = "productvalidator.ValidateCreateProductRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Description,
			validation.Required.Error("validation.required")),
		validation.Field(&req.IsSecondHand, validation.NotNil.Error("validation.required")),
		validation.Field(&req.ENTitle, validation.Required.Error("validation.required")),
		validation.Field(&req.BrandID, validation.Required.Error("validation.required"), validation.By(v.checkBrandIsExist)),
		validation.Field(&req.CategoryID, validation.Required.Error("validation.required"), validation.By(v.checkCategoryIsExist)),
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

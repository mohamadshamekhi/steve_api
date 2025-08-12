package productvalidator

import (
	"errors"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateProductRequest(req productparam.UpdateProductRequest) error {
	const op = "productvalidator.ValidateUpdateProductRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required.Error("validation.required")),
		validation.Field(&req.Name,
			validation.When(req.Name != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Description,
			validation.When(req.Description != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.IsSecondHand, validation.When(req.IsSecondHand != nil,
			validation.Required.Error("validation.required"))),
		validation.Field(&req.ENTitle, validation.When(req.ENTitle != nil,
			validation.Required.Error("validation.required"))),
		validation.Field(&req.BrandID,
			validation.When(req.BrandID != nil, validation.Required.Error("validation.required"),
				validation.By(v.checkBrandIsExist))),
		validation.Field(&req.CategoryID,
			validation.When(req.CategoryID != nil, validation.Required.Error("validation.required"),
				validation.By(v.checkCategoryIsExist))),
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

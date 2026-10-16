package productvalidator

import (
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
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
		validation.Field(&req.IsSecondHand, validation.NotNil.Error("validation.required")),
		validation.Field(&req.ENTitle, validation.When(req.ENTitle != nil,
			validation.Required.Error("validation.required"))),
		validation.Field(&req.BrandID,
			validation.When(req.BrandID != nil, validation.Required.Error("validation.required"),
				validation.By(v.checkBrandIsExist))),
		validation.Field(&req.CategoryID,
			validation.When(req.CategoryID != nil, validation.Required.Error("validation.required"),
				validation.By(v.checkCategoryIsExist))),
		validation.Field(&req.Price,
			validation.When(req.Price != nil,
				validation.Min(uint64(1)).Error("validation.invalid_format"))),
		validation.Field(&req.Kind,
			validation.When(req.Kind != nil, validation.By(validateKindPtr))),
		validation.Field(&req.SalesChannels, validation.By(validateSalesChannels)),
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

// validateKindPtr applies the same rule as validateKind to the update DTO's
// optional pointer field.
func validateKindPtr(value any) error {
	kind, ok := value.(*entity.ProductKind)
	if !ok || kind == nil || *kind == "" {
		return nil
	}
	if !kind.IsValid() {
		return errors.New("validation.invalid_format")
	}
	return nil
}

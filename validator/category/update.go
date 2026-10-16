package categoryvalidator

import (
	"errors"

	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateCategoryRequest(req categoryparam.UpdateCategoryRequest) error {
	const op = "categoryvalidator.ValidateUpdateCategoryRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.CategoryID,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Name,
			validation.When(req.Name != nil,
				validation.Required.Error("validation.required")),
		),
		validation.Field(&req.Description,
			validation.When(req.Description != nil,
				validation.Required.Error("validation.required"),
			),
		),
		validation.Field(&req.MediaID,
			validation.When(req.MediaID != nil,
				validation.Required.Error("validation.required"),
				validation.By(v.mediaIsExist),
			),
		),
		validation.Field(&req.ParentID,
			validation.When(req.ParentID != nil,
				validation.By(v.categoryIsExist),
			),
		),
		validation.Field(&req.AttributeSetID,
			validation.When(req.AttributeSetID != nil,
				validation.By(v.AttrIsExist),
			),
		),
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

package attributedefinitionvalidator

import (
	"errors"

	attributedefinitionparam "github.com/abolfazlnorzad/steve/param/attribute_definition"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateAttributeDefinitionRequest(req attributedefinitionparam.UpdateAttributeDefinitionRequest) error {
	const op = "attributedefinitionvalidator.ValidateUpdateAttributeDefinitionRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID,
			validation.Required.Error("validation.required"), validation.By(v.checkAttrIsExistByID)),
		validation.Field(&req.Unit, validation.When(req.Unit != nil, validation.Required.Error("validation.required"))),
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

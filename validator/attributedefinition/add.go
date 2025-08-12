package attributedefinitionvalidator

import (
	"errors"

	attributedefinitionparam "github.com/abolfazlnorzad/steve/param/attribute_definition"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateCreateAttributeDefinitionRequest(req attributedefinitionparam.AddAttributeDefinitionRequest) error {
	const op = "attributedefinitionvalidator.ValidateCreateAttributeDefinitionRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Key,
			validation.Required.Error("validation.required"), validation.By(v.checkKeyIsUnique))); err != nil {

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

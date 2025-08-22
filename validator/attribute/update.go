package attributevalidator

import (
	"errors"

	attributeparam "github.com/abolfazlnorzad/steve/param/attribute"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateAttributeRequest(req attributeparam.UpdateAttributeRequest) error {
	const op = "attributevalidator.ValidateCreateAttributeRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required.Error("validation.required"), validation.By(v.checkAttrIsExist)),
		validation.Field(&req.Name,
			validation.When(req.Name != nil,
				validation.Required.Error("validation.required")),
		),
		validation.Field(&req.AttributeDefinitionIDs,
			validation.When(req.AttributeDefinitionIDs != nil,
				validation.Required.Error("validation.required")),
		)); err != nil {

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

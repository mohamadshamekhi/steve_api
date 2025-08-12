package backofficevalidator

import (
	"errors"

	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateAddNewBackofficeUserRequest(req backofficeparam.AddNewBackofficeUserRequest) error {
	const op = "backofficevalidator.ValidateLoginWithPhoneNumberRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.PhoneNumber,
			validation.Required.Error("validation.required"),
			validation.By(v.isBackOfficeUserIsExist)),
		validation.Field(&req.Name,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Address,
			validation.Required.Error("validation.required")),
		validation.Field(&req.IsActive,
			validation.NotNil.Error("validation.required")),
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

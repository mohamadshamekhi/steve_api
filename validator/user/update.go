package uservalidator

import (
	"errors"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/constant"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

func (v Validator) ValidateStoreFrontUpdateProfileRequest(req userparam.StoreFrontUpdateProfileRequest) error {
	const op = "uservalidator.ValidateStoreFrontUpdateProfileRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required.Error("validation.required"), validation.By(v.isUserIsExist)),
		validation.Field(&req.PhoneNumber,
			validation.When(req.PhoneNumber != nil, validation.Required.Error("validation.required"),
				validation.Match(regexp.MustCompile(constant.PhoneNumberRegex)).Error("validation.phone_number_invalid"),
				validation.By(v.checkPhoneNumberIsUniqueForUpdate(req.UserID))),
			validation.When(req.Name != nil, validation.Required.Error("validation.required")),
		)); err != nil {

		var vErrs validation.Errors
		if errors.As(err, &vErrs) {
			fieldErrorKeys := make(map[string]string)
			for field, fieldErr := range vErrs {
				if fieldErr != nil {
					fieldErrorKeys[field] = fieldErr.Error()
				}
			}

			return richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("global.validation_failed").
				WithMeta(map[string]any{"fields": fieldErrorKeys})
		}

		return richerror.New(op).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong").
			WithWarpError(err)

	}

	//nolint
	return nil
}

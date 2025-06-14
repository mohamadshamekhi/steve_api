package backofficevalidatorauth

import (
	"context"
	"errors"
	backofficeauthparam "github.com/abolfazlnorzad/steve/param/backoffice/auth"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

func (v Validator) ValidateLoginWithPhoneNumberRequest(req backofficeauthparam.SendLoginCodeRequest) error {
	const op = "backofficevalidatorauth.ValidateLoginWithPhoneNumberRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.PhoneNumber,
			validation.Required.Error("validation.required"),
			validation.Match(regexp.MustCompile(phoneNumberRegex)).Error("validation.phone_number_invalid"),
			validation.By(v.isBackOfficeUserIsExist))); err != nil {

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

func (v Validator) isBackOfficeUserIsExist(value any) error {
	phoneNumber, ok := value.(string)
	if !ok {
		return richerror.New("validator.isBackOfficeUserIsExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	isExist, err := v.repo.IsBackOfficeUserIsExistByPhoneNumber(context.Background(), phoneNumber)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("auth.user_not_found")
	}

	return nil
}

func (v Validator) ValidateVerifyOTPCodeRequest(req backofficeauthparam.VerifyLoginCodeRequest) error {
	const op = "backofficevalidatorauth.ValidateVerifyOTPCodeRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Code, validation.Required.Error("validation.required")),
		validation.Field(&req.PhoneNumber,
			validation.Required.Error("validation.required"),
			validation.Match(regexp.MustCompile(phoneNumberRegex)).Error("validation.phone_number_invalid"),
			validation.By(v.isBackOfficeUserIsExist))); err != nil {

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

package uservalidator

import (
	"context"
	"errors"
	"regexp"

	userauthparam "github.com/abolfazlnorzad/steve/param/user/auth"
	"github.com/abolfazlnorzad/steve/pkg/constant"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUserLoginWithPhoneNumberRequest(req userauthparam.UserSendLoginCodeRequest) error {
	const op = "uservalidator.ValidateUserLoginWithPhoneNumberRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.PhoneNumber,
			validation.Required.Error("validation.required"),
			validation.Match(regexp.MustCompile(constant.PhoneNumberRegex)).Error("validation.phone_number_invalid"),
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
func (v Validator) isUserIsExistByPhoneNumber(value any) error {
	phoneNumber, ok := value.(string)
	if !ok {
		return richerror.New("validator.isUserIsExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	isExist, err := v.repo.CheckUserIsExistByPhoneNumber(context.Background(), phoneNumber)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("auth.user_not_found")
	}

	return nil
}

func (v Validator) isUserIsExistByID(value any) error {
	id, ok := value.(string)
	if !ok {
		return richerror.New("validator.isUserIsExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	isExist, err := v.repo.CheckUserIsExistByID(context.Background(), id)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("auth.user_not_found")
	}

	return nil
}

func (v Validator) ValidateUserVerifyOTPCodeRequest(req userauthparam.UserVerifyLoginCodeRequest) error {
	const op = "uservalidator.ValidateUserVerifyOTPCodeRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Code, validation.Required.Error("validation.required")),
		validation.Field(&req.PhoneNumber,
			validation.Required.Error("validation.required"),
			validation.Match(regexp.MustCompile(constant.PhoneNumberRegex)).Error("validation.phone_number_invalid"),
		)); err != nil {

		var vErrs validation.Errors
		if errors.As(err, &vErrs) {
			return richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("global.validation_failed").
				WithMeta(map[string]any{"fields": vErrs})
		}

		return richerror.New(op).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong").
			WithWarpError(err)

	}

	//nolint
	return nil
}

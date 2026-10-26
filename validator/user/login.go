package uservalidator

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	userauthparam "github.com/abolfazlnorzad/steve/param/user/auth"
	"github.com/abolfazlnorzad/steve/pkg/constant"
	"github.com/abolfazlnorzad/steve/pkg/phonenumber"
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
		return verifyValidationError(op, err)
	}

	// Name and birth date are mandatory only for a brand-new signup. Checked
	// here, before the OTP is consumed, so the user can resubmit the same code.
	exists, eErr := v.repo.CheckUserIsExistByPhoneNumber(context.Background(), phonenumber.Normalize(req.PhoneNumber))
	if eErr != nil || exists {
		// On lookup failure the service treats the user as existing too.
		return nil
	}

	now := time.Now()
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.By(func(any) error {
				if strings.TrimSpace(req.Name) == "" {
					return errors.New("validation.required")
				}
				return nil
			}),
			validation.RuneLength(2, 100).Error("validation.invalid_format"),
		),
		validation.Field(&req.BirthDate,
			validation.Required.Error("validation.required"),
			validation.Date(constant.BirthDateLayout).
				Min(now.AddDate(-120, 0, 0)).
				Max(now).
				Error("validation.birth_date_invalid").
				RangeError("validation.birth_date_invalid"),
		)); err != nil {
		return verifyValidationError(op, err)
	}

	return nil
}

func verifyValidationError(op string, err error) error {
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

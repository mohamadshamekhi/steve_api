package uservalidator

import (
	"context"
	"errors"
	"regexp"

	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/constant"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateStoreFrontUpdateProfileRequest(req userparam.StoreFrontUpdateProfileRequest) error {
	const op = "uservalidator.ValidateStoreFrontUpdateProfileRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required.Error("validation.required"), validation.By(v.isUserIsExistByID)),
		validation.Field(&req.PhoneNumber,
			validation.When(req.PhoneNumber != nil, validation.Required.Error("validation.required"),
				validation.Match(regexp.MustCompile(constant.PhoneNumberRegex)).Error("validation.phone_number_invalid"),
				validation.By(v.checkPhoneNumberIsUniqueForUpdate(req.UserID))),
			validation.When(req.Name != nil, validation.Required.Error("validation.required")),
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

func (v Validator) ValidateBackofficeUpdateUserRequest(req userparam.UpdateUserBackofficeRequest) error {
	const op = "uservalidator.ValidateBackofficeUpdateUserRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required.Error("validation.required"), validation.By(v.isUserIsExistByID)),
		validation.Field(&req.Name, validation.Required.Error("validation.required")),
	); err != nil {

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

func (v Validator) userIsExist(value any) error {
	id, ok := value.(string)
	if !ok {
		return richerror.New("uservalidator.userIsExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	resp, err := v.repo.CheckUserIsExistByID(context.Background(), id)
	if err != nil {
		return err
	}
	if resp {
		return nil
	}
	return errors.New("validation.user_is_already_exist")
}

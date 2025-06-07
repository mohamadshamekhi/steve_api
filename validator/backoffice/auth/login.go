package backofficevalidatorauth

import (
	"context"
	"errors"
	"fmt"
	backofficeauthparam "github.com/abolfazlnorzad/steve/param/backoffice/auth"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

func (v Validator) ValidateLoginWithPhoneNumberRequest(req backofficeauthparam.SendLoginCodeRequest) (map[string]string, error) {
	const op = "backofficevalidatorauth.ValidateLoginWithPhoneNumberRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.PhoneNumber,
			validation.Required,
			validation.Match(regexp.MustCompile(phoneNumberRegex)).Error(errmsg.ErrorMsgPhoneNumberIsNotValid),
			validation.By(v.isBackOfficeUserIsExist))); err != nil {
		fieldErrors := make(map[string]string)

		vErr := validation.Errors{}
		if errors.As(err, &vErr) {
			for key, value := range vErr {

				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}

		return fieldErrors, richerror.New(op).WithMessage(errmsg.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"req": req}).WithWarpError(err)
	}

	//nolint
	return nil, nil
}

func (v Validator) isBackOfficeUserIsExist(value any) error {
	phoneNumber, ok := value.(string)
	if !ok {
		return fmt.Errorf(errmsg.ErrorMsgSomethingWentWrong)
	}
	isExist, err := v.repo.IsBackOfficeUserIsExistByPhoneNumber(context.Background(), phoneNumber)
	if err != nil {
		return fmt.Errorf(errmsg.ErrorMsgSomethingWentWrong)
	}
	if !isExist {
		return fmt.Errorf(errmsg.ErrorMsgPhoneNumberOrPassIsIncorrect)
	}

	return nil
}

func (v Validator) ValidateVerifyOTPCodeRequest(req backofficeauthparam.VerifyLoginCodeRequest) (map[string]string, error) {
	const op = "backofficevalidatorauth.ValidateVerifyOTPCodeRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Code, validation.Required),
		validation.Field(&req.PhoneNumber,
			validation.Required,
			validation.Match(regexp.MustCompile(phoneNumberRegex)).Error(errmsg.ErrorMsgPhoneNumberIsNotValid),
			validation.By(v.isBackOfficeUserIsExist))); err != nil {
		fieldErrors := make(map[string]string)

		vErr := validation.Errors{}
		if errors.As(err, &vErr) {
			for key, value := range vErr {

				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}

		return fieldErrors, richerror.New(op).WithMessage(errmsg.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"req": req}).WithWarpError(err)
	}

	//nolint
	return nil, nil
}

package orderdeliveryvalidator

import (
	"errors"
	"regexp"

	orderdeliveryparam "github.com/abolfazlnorzad/steve/param/orderdelivery"
	"github.com/abolfazlnorzad/steve/pkg/constant"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Validator struct{}

func New() Validator {
	return Validator{}
}

func (v Validator) SendOTP(req orderdeliveryparam.SendOTPRequest) error {
	const op = "orderdeliveryvalidator.SendOTP"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Token, validation.Required.Error("validation.required")),
		validation.Field(&req.PhoneNumber, validation.Required.Error("validation.required"),
			validation.Match(regexp.MustCompile(constant.PhoneNumberRegex)).Error("validation.phone_number_invalid")),
	); err != nil {
		return fieldsErr(op, err)
	}
	return nil
}

func (v Validator) Verify(req orderdeliveryparam.VerifyRequest) error {
	const op = "orderdeliveryvalidator.Verify"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Token, validation.Required.Error("validation.required")),
		validation.Field(&req.PhoneNumber, validation.Required.Error("validation.required"),
			validation.Match(regexp.MustCompile(constant.PhoneNumberRegex)).Error("validation.phone_number_invalid")),
		validation.Field(&req.Code, validation.Required.Error("validation.required"),
			validation.Match(regexp.MustCompile(`^\d{5}$`)).Error("validation.incorrect_opt_code")),
	); err != nil {
		return fieldsErr(op, err)
	}
	return nil
}

func fieldsErr(op string, err error) error {
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

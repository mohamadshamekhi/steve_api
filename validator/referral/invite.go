package referralvalidator

import (
	"errors"
	"regexp"

	referralparam "github.com/abolfazlnorzad/steve/param/referral"
	"github.com/abolfazlnorzad/steve/pkg/constant"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateInviteByPhoneRequest(req referralparam.InviteByPhoneRequest) error {
	const op = "referralvalidator.ValidateInviteByPhoneRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.PhoneNumber,
			validation.Required.Error("validation.required"),
			validation.Match(regexp.MustCompile(constant.PhoneNumberRegex)).Error("validation.phone_number_invalid")),
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
	return nil
}

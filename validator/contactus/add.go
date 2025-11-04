package contactusvalidator

import (
	"errors"
	"regexp"

	contactusparam "github.com/abolfazlnorzad/steve/param/contactus"
	"github.com/abolfazlnorzad/steve/pkg/constant"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (v Validator) ValidateAddNewContactRequest(req contactusparam.AddNewContactRequest) error {
	const op = "contactusvalidator.ValidateAddNewContactRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Email,
			validation.Required.Error("validation.required"), is.Email),
		validation.Field(&req.Subject,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Description,
			validation.Required.Error("validation.required")),
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

package partnervalidator

import (
	"errors"
	"regexp"

	partnerparam "github.com/abolfazlnorzad/steve/param/partner"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var nationalIDRegex = regexp.MustCompile(`^\d{10}$`)

func (v Validator) ValidateApplyRequest(req partnerparam.ApplyRequest) error {
	const op = "partnervalidator.ValidateApplyRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.CompanyName,
			validation.Required.Error("validation.required"),
			validation.Length(2, 150).Error("validation.invalid_format")),
		validation.Field(&req.NationalID,
			validation.Required.Error("validation.required"),
			validation.Match(nationalIDRegex).Error("validation.invalid_format")),
		validation.Field(&req.BusinessField,
			validation.Required.Error("validation.required"),
			validation.Length(2, 150).Error("validation.invalid_format")),
		validation.Field(&req.City,
			validation.Required.Error("validation.required"),
			validation.Length(2, 100).Error("validation.invalid_format")),
	); err != nil {
		return wrap(op, err)
	}
	return nil
}

func wrap(op string, err error) error {
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

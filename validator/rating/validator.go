package ratingvalidator

import (
	"errors"

	ratingparam "github.com/abolfazlnorzad/steve/param/rating"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Validator struct {
}

func New() Validator {
	return Validator{}
}

func (v Validator) ValidateUpsertRatingRequest(req ratingparam.UpsertRatingRequest) error {
	const op = "backofficevalidatorauth.ValidateLoginWithPhoneNumberRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.UserID,
			validation.Required.Error("validation.required")),
		validation.Field(&req.ProductID,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Rating,
			validation.Required.Error("validation.required"), validation.Min(float64(1)), validation.Max(float64(5))),
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

	//nolint
	return nil
}

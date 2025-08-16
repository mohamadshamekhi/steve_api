package uservalidator

import (
	"errors"

	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateStoreFrontAddAddressRequest(req userparam.StoreFrontAddAddressRequest) error {
	const op = "uservalidator.ValidateStoreFrontAddAddressRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required.Error("validation.required"), validation.By(v.isUserIsExistByID)),
		validation.Field(&req.City, validation.Required.Error("validation.required")),
		validation.Field(&req.Title, validation.Required.Error("validation.required")),
		validation.Field(&req.Postcode, validation.Required.Error("validation.required")),
		validation.Field(&req.Province, validation.Required.Error("validation.required")),
		validation.Field(&req.Street, validation.Required.Error("validation.required"))); err != nil {

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

func (v Validator) ValidateStoreFrontUpdateAddressRequest(req userparam.StoreFrontUpdateAddressRequest) error {
	const op = "uservalidator.ValidateStoreFrontAddAddressRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required.Error("validation.required")),
		validation.Field(&req.UserID, validation.Required.Error("validation.required"), validation.By(v.isUserIsExistByID)),
		validation.Field(&req.City, validation.When(req.City != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Title, validation.When(req.City != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Postcode, validation.When(req.Postcode != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Province, validation.When(req.Province != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Street, validation.When(req.Street != nil, validation.Required.Error("validation.required")))); err != nil {

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

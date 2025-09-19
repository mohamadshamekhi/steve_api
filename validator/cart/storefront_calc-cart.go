package cartvalidator

import (
	"errors"

	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateCalcCartRequest(req cartparam.CalculateCartRequest) error {
	const op = "cartvalidator.ValidateCalcCartRequest"

	// todo - check discount code is exist and valid ?
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.DiscountCode, validation.Required.Error("validation.required")),
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

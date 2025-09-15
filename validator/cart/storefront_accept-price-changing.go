package cartvalidator

import (
	"errors"

	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateAcceptPriceChangingRequest(req cartparam.AcceptPriceChangeRequest) error {
	const op = "cartvalidator.ValidateAcceptPriceChangingRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.VariantID, validation.Required.Error("validation.required")),
	); err != nil {

		// todo - check variant id is exist ?
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

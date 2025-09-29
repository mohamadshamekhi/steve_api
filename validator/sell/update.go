package sellvalidator

import (
	"errors"

	sellparam "github.com/abolfazlnorzad/steve/param/sell"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateSellRequest(req sellparam.UpdateSellRequest) error {
	const op = "sellvalidator.UpdateSellRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID,
			validation.Required.Error("validation.required")),
		// todo -check brand and product and media is exist ?
		validation.Field(&req.Status,
			validation.Required.Error("validation.required"),
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

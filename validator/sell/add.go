package sellvalidator

import (
	"errors"

	sellparam "github.com/abolfazlnorzad/steve/param/sell"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateAddSellRequest(req sellparam.AddSellRequest) error {
	const op = "sellvalidator.ValidateAddSellRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.UserID,
			validation.Required.Error("validation.required")),
		// todo -check brand and product is exist ?
		validation.Field(&req.BrandID,
			validation.Required.Error("validation.required")),
		validation.Field(&req.ProductID,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Description,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Color,
			validation.Required.Error("validation.required")),
		validation.Field(&req.PhysicalCondition,
			validation.Required.Error("validation.required")),
		validation.Field(&req.BatteryHealth,
			validation.Required.Error("validation.required")),
		validation.Field(&req.RegisterStatus,
			validation.Required.Error("validation.required")),
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

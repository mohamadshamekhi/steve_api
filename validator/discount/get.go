package discountvalidator

import (
	"errors"

	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateGetDiscountRequest(req discountparam.GetSingleDiscountRequest) error {
	const op = "discountvalidator.ValidateGetDiscountRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required.Error("validation.required"), validation.By(v.checkDiscountIsExist)),
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

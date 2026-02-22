package cartvalidator

import (
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) UpdateStatusCartValidator(req cartparam.UpdateCartStatusRequest) error {
	const op = "cartvalidator.UpdateStatusCartValidator"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Status,
			validation.Required.Error("validation.required"),
			validation.In(
				entity.CartActive,
				entity.CartInactive,
				entity.CartPaid,
				entity.CartCanceled,
				entity.PendingPayment,
			).Error("validation.invalid_status"),
		),
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

package ordervalidator

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) UpdateStatusOrderValidator(ctx context.Context, req orderparam.UpdateStatusOrderRequest) error {
	const op = "ordervalidator.UpdateStatusOrderValidator"
	if err := validation.ValidateStructWithContext(ctx, &req,
		validation.Field(&req.Status,
			validation.Required.Error("validation.required"),
			validation.In(
				entity.OrderStatusPendingPayment,
				entity.OrderStatusProcessing,
				entity.OrderStatusShipped,
				entity.OrderStatusDelivered,
				entity.OrderStatusCancelled,
				entity.OrderStatusRefunded,
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

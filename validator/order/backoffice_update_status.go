package ordervalidator

import (
	"context"
	"errors"
	"regexp"
	"strings"

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
		validation.Field(&req.TrackingCode,
			validation.When(req.Status == entity.OrderStatusShipped,
				validation.By(requiredTrimmed),
				validation.Length(4, 64).Error("validation.invalid_format"),
				validation.Match(trackingCodeRegex).Error("validation.invalid_format"),
			),
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

// trackingCodeRegex accepts courier tracking numbers: Post's all-digit codes
// as well as alphanumeric ones with dashes (Tipax, Chapar, …).
var trackingCodeRegex = regexp.MustCompile(`^\s*[A-Za-z0-9-]+\s*$`)

func requiredTrimmed(value interface{}) error {
	s, _ := value.(string)
	if strings.TrimSpace(s) == "" {
		return errors.New("validation.required")
	}
	return nil
}

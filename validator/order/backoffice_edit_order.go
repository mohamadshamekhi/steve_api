package ordervalidator

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) EditOrderValidator(ctx context.Context, req orderparam.EditOrderRequest) error {
	const op = "ordervalidator.EditOrderValidator"

	err := validation.ValidateStructWithContext(ctx, &req,
		validation.Field(&req.Reason,
			validation.Required.Error("validation.required"),
			validation.Length(3, 500).Error("validation.invalid_input"),
		),
		validation.Field(&req.PaymentMethod,
			validation.When(req.PaymentMethod != "",
				validation.By(func(any) error {
					if !req.PaymentMethod.IsValid() {
						return errors.New("validation.invalid_input")
					}
					return nil
				}),
			),
		),
		validation.Field(&req.Items, validation.Each(validation.By(func(value any) error {
			it, _ := value.(orderparam.EditOrderItem)
			if it.Quantity == 0 {
				return errors.New("validation.invalid_input")
			}
			if it.IsCustom {
				if it.ProductName == "" || it.UnitPrice == nil {
					return errors.New("validation.required")
				}
				if it.Kind != "" && it.Kind != entity.ProductKindPhysical && it.Kind != entity.ProductKindService {
					return errors.New("validation.invalid_input")
				}
				return nil
			}
			if it.VariantID == "" {
				return errors.New("validation.required")
			}
			return nil
		}))),
	)
	if err == nil {
		return nil
	}
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

package ordervalidator

import (
	"context"
	"errors"

	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) RegisterOrderValidator(ctx context.Context, req orderparam.RegisterOrderRequest) error {
	const op = "ordervalidator.RegisterOrderValidator"

	if err := validation.ValidateStructWithContext(ctx, &req,
		validation.Field(&req.DiscountCode, validation.Required.Error("validation.required"), validation.By(v.checkDiscountIsValid)),
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

func (v Validator) checkDiscountIsValid(value any) error {

	const op = "ordervalidator.checkDiscountIsValid"
	code, _ := value.(string)
	_, err := v.discountProvider.GetByCode(context.Background(), discountparam.GetSingleDiscountByCodeRequest{Code: code})
	if err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	return nil
}

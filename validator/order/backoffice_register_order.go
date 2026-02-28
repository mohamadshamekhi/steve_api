package ordervalidator

import (
	"context"
	"errors"

	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) RegisterOrderBackofficeValidator(ctx context.Context, req orderparam.RegisterOrderBackofficeRequest) error {
	const op = "ordervalidator.RegisterOrderBackofficeValidator"
	if err := validation.ValidateStructWithContext(ctx, &req,
		validation.Field(&req.UserID, validation.Required.Error("validation.required")),
		validation.Field(&req.DiscountCode,
			validation.When(req.DiscountCode != "",
				validation.By(v.discountCodeIsValid(ctx, req.UserID)),
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

	//nolint
	return nil
}

func (v Validator) discountCodeIsValid(ctx context.Context, userID string) validation.RuleFunc {
	return func(value any) error {
		var code string
		switch v := value.(type) {
		case string:
			code = v
		case *string:
			if v == nil {
				return nil
			}
			code = *v
		default:
			return richerror.New("brandvalidator.BrandIsUnique.type_assertion").
				WithKind(richerror.KindUnexpected).
				WithWarpError(errors.New("unexpected type for brand name"))
		}

		if code == "" {
			return nil
		}
		cart, err := v.cartRepo.GetUserActiveCart(ctx, userID)
		if err != nil {
			return err
		}
		resp, err := v.discountProvider.GetValidDiscountsForCart(ctx, discountparam.GetValidDiscountsForSpecificCartRequest{DiscountCode: code, Cart: *cart})
		if err != nil {
			return err
		}

		for _, d := range resp.Discounts {
			if d.Code == code {
				return nil
			}
		}
		return richerror.New("validator.discountCodeIsValid").
			WithKind(richerror.KindInvalid).
			WithMessageKey("discount.code_is_not_valid")
	}
}

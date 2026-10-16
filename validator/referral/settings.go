package referralvalidator

import (
	"errors"

	referralparam "github.com/abolfazlnorzad/steve/param/referral"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var discountTypes = []any{"percentage", "fixed_amount"}

func (v Validator) ValidateUpdateReferralSettingsRequest(req referralparam.UpdateReferralSettingsRequest) error {
	const op = "referralvalidator.ValidateUpdateReferralSettingsRequest"

	rules := []*validation.FieldRules{
		validation.Field(&req.RewardThreshold,
			validation.Required.Error("validation.required"),
			validation.Min(1).Error("validation.invalid_format")),
		validation.Field(&req.RewardDiscountType,
			validation.Required.Error("validation.required"),
			validation.In(discountTypes...).Error("validation.invalid_format")),
		validation.Field(&req.RewardDiscountValue,
			validation.Required.Error("validation.required"),
			validation.By(percentRule(req.RewardDiscountType, req.RewardDiscountValue))),
		validation.Field(&req.RewardValidDays,
			validation.Required.Error("validation.required"),
			validation.Min(1).Error("validation.invalid_format")),
		validation.Field(&req.WelcomeDiscountType,
			validation.When(req.WelcomeDiscountEnabled,
				validation.Required.Error("validation.required"),
				validation.In(discountTypes...).Error("validation.invalid_format"))),
		validation.Field(&req.WelcomeDiscountValue,
			validation.When(req.WelcomeDiscountEnabled,
				validation.Required.Error("validation.required"),
				validation.By(percentRule(req.WelcomeDiscountType, req.WelcomeDiscountValue)))),
		validation.Field(&req.WelcomeValidDays,
			validation.When(req.WelcomeDiscountEnabled,
				validation.Required.Error("validation.required"),
				validation.Min(1).Error("validation.invalid_format"))),
	}

	if err := validation.ValidateStruct(&req, rules...); err != nil {
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

// percentRule caps a percentage discount at 100; fixed-amount values are left
// unbounded (any Toman amount is legitimate).
func percentRule(discountType string, _ uint) validation.RuleFunc {
	return func(value any) error {
		if discountType != "percentage" {
			return nil
		}
		v, ok := value.(uint)
		if !ok || v == 0 || v > 100 {
			return errors.New("validation.invalid_format")
		}
		return nil
	}
}

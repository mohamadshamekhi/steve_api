package birthdayvalidator

import (
	"errors"

	birthdayparam "github.com/abolfazlnorzad/steve/param/birthday"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var discountTypes = []any{"percentage", "fixed_amount"}

func (v Validator) ValidateUpdateSettingsRequest(req birthdayparam.UpdateSettingsRequest) error {
	const op = "birthdayvalidator.ValidateUpdateSettingsRequest"

	if err := validation.ValidateStruct(&req,
		// days_before may legitimately be 0 (send on the birthday), so no Required.
		validation.Field(&req.DaysBefore,
			validation.Min(0).Error("validation.invalid_format"),
			validation.Max(30).Error("validation.invalid_format")),
		validation.Field(&req.SendHour,
			validation.Min(0).Error("validation.invalid_format"),
			validation.Max(23).Error("validation.invalid_format")),
		validation.Field(&req.DiscountType,
			validation.Required.Error("validation.required"),
			validation.In(discountTypes...).Error("validation.invalid_format")),
		validation.Field(&req.DiscountValue,
			validation.Required.Error("validation.required"),
			validation.By(percentRule(req.DiscountType))),
		validation.Field(&req.ValidDays,
			validation.Required.Error("validation.required"),
			validation.Min(1).Error("validation.invalid_format"),
			validation.Max(365).Error("validation.invalid_format")),
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

// percentRule caps a percentage discount at 100; fixed-amount values are left
// unbounded (any Toman amount is legitimate).
func percentRule(discountType string) validation.RuleFunc {
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

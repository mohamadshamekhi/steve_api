package discountvalidator

import (
	"errors"
	"time"

	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateDiscountRequest(req discountparam.UpdateDiscountRequest) error {
	const op = "discountvalidator.ValidateUpdateDiscountRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required.Error("validation.required"), validation.By(v.checkDiscountIsExist)),
		validation.Field(&req.Value, validation.Required.Error("validation.required")),
		validation.Field(&req.ExpiresAt, validation.Required.Error("validation.required"), validation.Date(time.RFC3339).Error("validation.invalid_date_format")),
		validation.Field(&req.MinimumPurchaseAmount,
			validation.When(req.MinimumPurchaseAmount != nil,
				validation.Min(0).Error("validation.must_be_positive_number"),
			),
		),
		validation.Field(&req.ApplicableCategoryIDs, validation.By(v.categoriesExist)),
		validation.Field(&req.ApplicableProductIDs, validation.By(v.productsExist)),
		validation.Field(&req.ApplicableUserIDs, validation.By(v.usersExist)),
		validation.Field(&req.TotalUsageLimit,
			validation.When(req.MinimumPurchaseAmount != nil,
				validation.Min(0).Error("validation.must_be_positive_number"),
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

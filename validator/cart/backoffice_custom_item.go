package cartvalidator

import (
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateSetItemUnitPriceBackofficeRequest(req cartparam.SetItemUnitPriceBackofficeRequest) error {
	const op = "cartvalidator.ValidateSetItemUnitPriceBackofficeRequest"
	if req.UnitPrice < 0 {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("cart.unit_price_invalid")
	}
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required.Error("validation.required"), validation.By(v.userIsExistByID)),
		validation.Field(&req.VariantID, validation.Required.Error("validation.required")),
	); err != nil {
		return wrapValidationError(op, err)
	}
	return nil
}

// ValidateUpsertCustomItemBackofficeRequest checks a hand-typed invoice line.
// There is no catalog entry to fall back on, so the title and the price the
// salesperson typed are all the invoice will ever have.
func (v Validator) ValidateUpsertCustomItemBackofficeRequest(req cartparam.UpsertCustomItemBackofficeRequest) error {
	const op = "cartvalidator.ValidateUpsertCustomItemBackofficeRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required.Error("validation.required"), validation.By(v.userIsExistByID)),
		validation.Field(&req.Title, validation.Required.Error("validation.required"), validation.Length(2, 200)),
		validation.Field(&req.UnitPrice, validation.Min(int64(0)).Error("validation.min")),
		validation.Field(&req.Quantity, validation.Required.Error("validation.required"), validation.Min(1).Error("validation.min")),
		validation.Field(&req.Kind, validation.By(kindIsValid)),
	); err != nil {
		return wrapValidationError(op, err)
	}

	//nolint
	return nil
}

// ValidateSetItemManualDiscountBackofficeRequest checks a per-line counter
// discount. The amount is only capped against the line at registration time,
// where the line's real price is known.
func (v Validator) ValidateSetItemManualDiscountBackofficeRequest(req cartparam.SetItemManualDiscountBackofficeRequest) error {
	const op = "cartvalidator.ValidateSetItemManualDiscountBackofficeRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required.Error("validation.required"), validation.By(v.userIsExistByID)),
		validation.Field(&req.VariantID, validation.Required.Error("validation.required")),
	); err != nil {
		return wrapValidationError(op, err)
	}

	//nolint
	return nil
}

// ValidateSetCartManualDiscountBackofficeRequest checks an invoice-level
// counter discount.
func (v Validator) ValidateSetCartManualDiscountBackofficeRequest(req cartparam.SetCartManualDiscountBackofficeRequest) error {
	const op = "cartvalidator.ValidateSetCartManualDiscountBackofficeRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required.Error("validation.required"), validation.By(v.userIsExistByID)),
	); err != nil {
		return wrapValidationError(op, err)
	}

	//nolint
	return nil
}

// kindIsValid accepts an empty kind, which is read as a physical good.
func kindIsValid(value any) error {
	const op = "cartvalidator.kindIsValid"

	kind, ok := value.(entity.ProductKind)
	if !ok || kind == "" {
		return nil
	}
	if !kind.IsValid() {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("product.invalid_kind")
	}
	return nil
}

// wrapValidationError turns ozzo's field errors into the richerror shape the
// central error handler knows how to translate.
func wrapValidationError(op string, err error) error {
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

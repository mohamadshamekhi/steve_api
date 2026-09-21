package purchasevalidator

import (
	"errors"

	purchaseparam "github.com/abolfazlnorzad/steve/param/purchase"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Validator struct{}

func New() Validator {
	return Validator{}
}

// validateItem is shared by create and update: an invoice line must name a
// real variant, move at least one unit, and carry a cost — a zero cost would
// poison the weighted average with free goods.
func validateItem(item purchaseparam.CreatePurchaseItemRequest) error {
	return validation.ValidateStruct(&item,
		validation.Field(&item.ProductID, validation.Required.Error("validation.required")),
		validation.Field(&item.VariantID, validation.Required.Error("validation.required")),
		validation.Field(&item.Quantity, validation.Min(uint64(1)).Error("validation.invalid_format")),
		validation.Field(&item.UnitCost, validation.Min(uint64(1)).Error("validation.invalid_format")),
	)
}

func (v Validator) ValidateCreateRequest(req purchaseparam.CreatePurchaseRequest) error {
	const op = "purchasevalidator.ValidateCreateRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.SupplierID, validation.Required.Error("validation.required")),
		validation.Field(&req.Items,
			validation.Required.Error("validation.required"),
			validation.Length(1, 0).Error("validation.invalid_format")),
		validation.Field(&req.Note, validation.Length(0, 1000).Error("validation.invalid_format")),
	); err != nil {
		return wrap(op, err)
	}

	for _, item := range req.Items {
		if err := validateItem(item); err != nil {
			return wrap(op, err)
		}
	}
	return nil
}

func (v Validator) ValidateUpdateRequest(req purchaseparam.UpdatePurchaseRequest) error {
	const op = "purchasevalidator.ValidateUpdateRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.SupplierID,
			validation.When(req.SupplierID != nil,
				validation.Required.Error("validation.required"))),
		validation.Field(&req.Note,
			validation.When(req.Note != nil,
				validation.Length(0, 1000).Error("validation.invalid_format"))),
	); err != nil {
		return wrap(op, err)
	}

	for _, item := range req.Items {
		if err := validateItem(item); err != nil {
			return wrap(op, err)
		}
	}
	return nil
}

func wrap(op string, err error) error {
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

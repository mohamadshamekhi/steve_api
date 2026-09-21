package inventoryvalidator

import (
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Validator struct{}

func New() Validator {
	return Validator{}
}

func (v Validator) ValidateAdjustRequest(req inventoryparam.AdjustStockRequest) error {
	const op = "inventoryvalidator.ValidateAdjustRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ProductID, validation.Required.Error("validation.required")),
		validation.Field(&req.VariantID, validation.Required.Error("validation.required")),
		validation.Field(&req.Note, validation.Length(0, 500).Error("validation.invalid_format")),
	); err != nil {
		return wrap(op, err)
	}

	// Required.Error() cannot express "non-zero, either sign", so the one
	// rule that matters here is checked directly.
	if req.Quantity == 0 {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("inventory.adjustment_quantity_required")
	}
	return nil
}

func (v Validator) ValidateSetStockConfigRequest(req inventoryparam.SetStockConfigRequest) error {
	const op = "inventoryvalidator.ValidateSetStockConfigRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ProductID, validation.Required.Error("validation.required")),
		validation.Field(&req.VariantID, validation.Required.Error("validation.required")),
	); err != nil {
		return wrap(op, err)
	}

	if req.StockMode != nil && !req.StockMode.IsValid() {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("inventory.invalid_stock_mode").
			WithMeta(map[string]any{"allowed": []entity.StockMode{
				entity.StockModeTracked, entity.StockModeOnDemand, entity.StockModeNone,
			}})
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

package productvalidator

import (
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateCreateProductRequest(req productparam.CreateProductRequest) error {
	const op = "productvalidator.ValidateCreateProductRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Description,
			validation.Required.Error("validation.required")),
		validation.Field(&req.IsSecondHand, validation.NotNil.Error("validation.required")),
		validation.Field(&req.ENTitle, validation.Required.Error("validation.required")),
		// A service (repair, Apple-ID, ...) has no manufacturer and no product
		// category, so brand and category are required for goods only. Each
		// guard tests the same condition — the item's kind — rather than
		// borrowing a neighbouring field's check.
		validation.Field(&req.BrandID,
			validation.When(req.Kind != entity.ProductKindService,
				validation.Required.Error("validation.required")),
			validation.When(req.BrandID != "", validation.By(v.checkBrandIsExist))),
		validation.Field(&req.CategoryID,
			validation.When(req.Kind != entity.ProductKindService,
				validation.Required.Error("validation.required")),
			validation.When(req.CategoryID != "", validation.By(v.checkCategoryIsExist))),
		// A service is created complete: it has no variant step afterwards, so
		// its fee has to arrive with it. Goods price per variant instead.
		validation.Field(&req.Price,
			validation.When(req.Kind == entity.ProductKindService,
				validation.Min(uint64(1)).Error("validation.required"))),
		validation.Field(&req.Kind, validation.By(validateKind)),
		validation.Field(&req.SalesChannels, validation.By(validateSalesChannels)),
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

// validateKind accepts an empty kind (the entity defaults it to physical) and
// otherwise requires one of the two known values, so a typo can never be
// silently persisted as an unknown kind.
func validateKind(value any) error {
	kind, ok := value.(entity.ProductKind)
	if !ok || kind == "" {
		return nil
	}
	if !kind.IsValid() {
		return errors.New("validation.invalid_format")
	}
	return nil
}

// validateSalesChannels rejects unknown channels and a list that names none of
// them — an item sellable nowhere is never what the caller meant.
func validateSalesChannels(value any) error {
	channels, ok := value.([]entity.SalesChannel)
	if !ok || len(channels) == 0 {
		return nil
	}
	for _, ch := range channels {
		if !ch.IsValid() {
			return errors.New("validation.invalid_format")
		}
	}
	return nil
}

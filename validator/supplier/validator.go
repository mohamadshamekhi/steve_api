package suppliervalidator

import (
	"errors"
	"regexp"

	supplierparam "github.com/abolfazlnorzad/steve/param/supplier"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var phoneRegex = regexp.MustCompile(`^09\d{9}$`)

type Validator struct{}

func New() Validator {
	return Validator{}
}

func (v Validator) ValidateCreateRequest(req supplierparam.CreateSupplierRequest) error {
	const op = "suppliervalidator.ValidateCreateRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.Required.Error("validation.required"),
			validation.Length(2, 150).Error("validation.invalid_format")),
		validation.Field(&req.ShopName,
			validation.Length(0, 150).Error("validation.invalid_format")),
		validation.Field(&req.PhoneNumber,
			validation.Match(phoneRegex).Error("validation.invalid_format")),
		validation.Field(&req.Address,
			validation.Length(0, 500).Error("validation.invalid_format")),
		validation.Field(&req.Note,
			validation.Length(0, 1000).Error("validation.invalid_format")),
	); err != nil {
		return wrap(op, err)
	}
	return nil
}

func (v Validator) ValidateUpdateRequest(req supplierparam.UpdateSupplierRequest) error {
	const op = "suppliervalidator.ValidateUpdateRequest"

	// Each optional field guards on itself: a shared When() would silently
	// disable a rule whenever the other field happens to be absent.
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.When(req.Name != nil,
				validation.Required.Error("validation.required"),
				validation.Length(2, 150).Error("validation.invalid_format"))),
		validation.Field(&req.ShopName,
			validation.When(req.ShopName != nil,
				validation.Length(0, 150).Error("validation.invalid_format"))),
		validation.Field(&req.PhoneNumber,
			validation.When(req.PhoneNumber != nil,
				validation.Match(phoneRegex).Error("validation.invalid_format"))),
		validation.Field(&req.Address,
			validation.When(req.Address != nil,
				validation.Length(0, 500).Error("validation.invalid_format"))),
		validation.Field(&req.Note,
			validation.When(req.Note != nil,
				validation.Length(0, 1000).Error("validation.invalid_format"))),
	); err != nil {
		return wrap(op, err)
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

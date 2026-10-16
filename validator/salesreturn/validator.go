package salesreturnvalidator

import (
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	salesreturnparam "github.com/abolfazlnorzad/steve/param/salesreturn"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Validator struct{}

func New() Validator {
	return Validator{}
}

func (v Validator) ValidateCreateRequest(req salesreturnparam.CreateSalesReturnRequest) error {
	const op = "salesreturnvalidator.ValidateCreateRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.OrderID,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Items,
			validation.Required.Error("validation.required"),
			validation.Each(validation.By(validateItem))),
		validation.Field(&req.RefundMethod, validation.By(validateMethod)),
		validation.Field(&req.Reason,
			validation.Length(0, 500).Error("validation.invalid_format")),
		validation.Field(&req.Note,
			validation.Length(0, 1000).Error("validation.invalid_format")),
	); err != nil {
		return wrap(op, err)
	}
	return nil
}

func (v Validator) ValidateMarkRefundedRequest(req salesreturnparam.MarkRefundedRequest) error {
	const op = "salesreturnvalidator.ValidateMarkRefundedRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.RefundMethod, validation.By(validateMethod)),
	); err != nil {
		return wrap(op, err)
	}
	return nil
}

func validateItem(value any) error {
	item, ok := value.(salesreturnparam.CreateSalesReturnItem)
	if !ok {
		return errors.New("validation.invalid_value")
	}
	if item.VariantID == "" {
		return errors.New("validation.required")
	}
	if item.Quantity == 0 {
		return errors.New("validation.invalid_value")
	}
	return nil
}

func validateMethod(value any) error {
	m, ok := value.(entity.PaymentMethod)
	if !ok || m == "" {
		return nil
	}
	if !m.IsValid() {
		return errors.New("validation.invalid_value")
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

package accountingvalidator

import (
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Validator struct{}

func New() Validator {
	return Validator{}
}

// validCategory rejects an unknown expense category, so the breakdown report
// can never grow a bucket nobody defined.
var validCategory = validation.By(func(value any) error {
	c, ok := value.(entity.ExpenseCategory)
	if !ok || c == "" {
		return nil
	}
	if !c.IsValid() {
		return errors.New("validation.invalid_value")
	}
	return nil
})

var validMethod = validation.By(func(value any) error {
	m, ok := value.(entity.PaymentMethod)
	if !ok || m == "" {
		return nil
	}
	if !m.IsValid() {
		return errors.New("validation.invalid_value")
	}
	return nil
})

func (v Validator) ValidateCreateExpenseRequest(req accountingparam.CreateExpenseRequest) error {
	const op = "accountingvalidator.ValidateCreateExpenseRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Title,
			validation.Required.Error("validation.required"),
			validation.Length(2, 200).Error("validation.invalid_format")),
		validation.Field(&req.Category, validCategory),
		validation.Field(&req.Amount,
			validation.Required.Error("validation.required")),
		validation.Field(&req.PaymentMethod, validMethod),
		validation.Field(&req.Note,
			validation.Length(0, 1000).Error("validation.invalid_format")),
	); err != nil {
		return wrap(op, err)
	}
	return nil
}

func (v Validator) ValidateUpdateExpenseRequest(req accountingparam.UpdateExpenseRequest) error {
	const op = "accountingvalidator.ValidateUpdateExpenseRequest"

	// Each optional field guards on itself; a shared When() would silently
	// disable a rule whenever the other field happens to be absent.
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Title,
			validation.When(req.Title != nil,
				validation.Required.Error("validation.required"),
				validation.Length(2, 200).Error("validation.invalid_format"))),
		validation.Field(&req.Category,
			validation.When(req.Category != nil, validation.By(func(value any) error {
				if req.Category == nil || req.Category.IsValid() {
					return nil
				}
				return errors.New("validation.invalid_value")
			}))),
		validation.Field(&req.Amount,
			validation.When(req.Amount != nil,
				validation.Required.Error("validation.required"))),
		validation.Field(&req.PaymentMethod,
			validation.When(req.PaymentMethod != nil, validation.By(func(value any) error {
				if req.PaymentMethod == nil || req.PaymentMethod.IsValid() {
					return nil
				}
				return errors.New("validation.invalid_value")
			}))),
		validation.Field(&req.Note,
			validation.When(req.Note != nil,
				validation.Length(0, 1000).Error("validation.invalid_format"))),
	); err != nil {
		return wrap(op, err)
	}
	return nil
}

func (v Validator) ValidateCreateCashTransactionRequest(req accountingparam.CreateCashTransactionRequest) error {
	const op = "accountingvalidator.ValidateCreateCashTransactionRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Direction,
			validation.Required.Error("validation.required"),
			validation.In(entity.CashIn, entity.CashOut).Error("validation.invalid_value")),
		validation.Field(&req.Method, validMethod),
		validation.Field(&req.Amount,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Description,
			validation.Length(0, 500).Error("validation.invalid_format")),
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

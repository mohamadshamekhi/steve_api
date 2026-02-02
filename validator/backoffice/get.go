package backofficevalidator

import (
	"context"
	"errors"

	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateGetBackofficeRequest(req backofficeparam.GetBackofficeRequest) error {
	const op = "backofficevalidator.ValidateGetBackofficeRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID,
			validation.Required.Error("validation.required"),
			validation.By(v.backofficeIsExist)),
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

func (v Validator) backofficeIsExist(value any) error {
	id, ok := value.(string)
	if !ok {
		return richerror.New("backofficevalidator.backofficeIsExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	resp, err := v.repo.IsBackOfficeUserIsExistByID(context.Background(), id)
	if err != nil {
		return err
	}
	if resp {
		return nil
	}
	return errors.New("validation.backoffice_is_not_exist")
}

package brandvalidator

import (
	"context"
	"errors"

	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateGetBrandRequest(req brandparam.GetBrandRequest) error {
	const op = "backofficevalidatorauth.ValidateGetBrandRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID,
			validation.Required.Error("validation.required"),
			validation.By(v.brandIsExist)),
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

func (v Validator) brandIsExist(value any) error {
	id, ok := value.(string)
	if !ok {
		return richerror.New("brandvalidator.BrandIsUnique").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	resp, err := v.repo.CheckBrandIDIsExist(context.Background(), id)
	if err != nil {
		return err
	}
	if resp {
		return nil
	}
	return errors.New("validation.brand_is_already_exist")
}

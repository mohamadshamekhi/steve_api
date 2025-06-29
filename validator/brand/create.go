package brandvalidator

import (
	"context"
	"errors"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateCreateBrandRequest(req brandparam.CreateBrandRequest) error {
	const op = "backofficevalidatorauth.ValidateLoginWithPhoneNumberRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.Required.Error("validation.required"),
			validation.By(v.brandIsExist))); err != nil {

		var vErrs validation.Errors
		if errors.As(err, &vErrs) {
			fieldErrorKeys := make(map[string]string)
			for field, errMessage := range vErrs {
				if errMessage != nil {
					fieldErrorKeys[field] = errMessage.Error()
				}
			}

			return richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("global.validation_failed").
				WithMeta(map[string]any{"fields": fieldErrorKeys})
		}

		return richerror.New(op).WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong").WithWarpError(err)

	}

	//nolint
	return nil
}

func (v Validator) brandIsExist(value any) error {
	name, ok := value.(string)
	if !ok {
		return richerror.New("brandvalidator.BrandIsUnique").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	resp, err := v.repo.CheckBrandNameIsExist(context.Background(), name)
	if err != nil {
		return err
	}
	if resp {
		return richerror.New("brandvalidator.BrandIsExist").WithKind(richerror.KindBadRequest).WithMessageKey("validation.brand_is_already_exist")
	}
	return nil
}

package backofficevalidator

import (
	"context"
	"errors"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateBackofficeUserRequest(req backofficeparam.UpdateBackofficeUserRequest) error {
	const op = "backofficevalidator.ValidateUpdateBackofficeUserRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required.Error("validation.required"), validation.By(v.checkBackofficeUserExist)),
		validation.Field(&req.PhoneNumber,
			validation.When(req.PhoneNumber != nil, validation.Required.Error("validation.required"),
				validation.By(v.checkPhoneNumberIsUniqueForUpdate(req.ID)))),
		validation.Field(&req.Name,
			validation.When(req.Name != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Address,
			validation.When(req.Address != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.IsActive,
			validation.When(req.IsActive != nil, validation.NotNil.Error("validation.required"))),
	); err != nil {

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

func (v Validator) checkBackofficeUserExist(value any) error {
	id, ok := value.(string)
	if !ok {
		return richerror.New("validator.isBackOfficeUserIsExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	_, err := v.repo.GetBackOfficeUserByID(context.Background(), id)
	if err != nil {
		var rErr *richerror.RichError
		if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
			return errors.New(rErr.GetMessageKey())
		}
		return err
	}

	return nil
}

func (v Validator) checkPhoneNumberIsUniqueForUpdate(currentBackofficeUserID string) validation.RuleFunc {
	return func(value interface{}) error {
		var phonenumber string
		switch v := value.(type) {
		case string:
			phonenumber = v
		case *string:
			if v == nil {
				return nil
			}
			phonenumber = *v
		default:
			return richerror.New("brandvalidator.BrandIsUnique.type_assertion").
				WithKind(richerror.KindUnexpected).
				WithWarpError(errors.New("unexpected type for brand name"))
		}

		existingBrand, err := v.repo.GetBackOfficeUserByPhoneNumber(context.Background(), phonenumber)
		if err != nil {
			var rErr *richerror.RichError
			if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
				return nil
			}
			return err
		}

		if existingBrand.ID != currentBackofficeUserID {
			return errors.New("validation.phone_number_is_already_exist")
		}

		return nil
	}
}

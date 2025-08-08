package brandvalidator

import (
	"context"
	"errors"

	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateCreateBrandRequest(req brandparam.CreateBrandRequest) error {
	const op = "backofficevalidatorauth.ValidateLoginWithPhoneNumberRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.Required.Error("validation.required"),
			validation.By(v.brandIsUnique)),
		validation.Field(&req.MediaID,
			validation.Required.Error("validation.required"),
			validation.By(v.mediaBrandIsExist),
		)); err != nil {

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

func (v Validator) brandIsUnique(value any) error {
	var name string
	switch v := value.(type) {
	case string:
		name = v
	case *string:
		if v == nil {
			return nil
		}
		name = *v
	default:
		return richerror.New("brandvalidator.BrandIsUnique.type_assertion").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type for brand name"))
	}

	if name == "" {
		return nil
	}
	resp, err := v.repo.CheckBrandNameIsExist(context.Background(), name)
	if err != nil {
		return err
	}
	if resp {
		return errors.New("validation.brand_is_already_exist")
	}
	return nil
}

func (v Validator) mediaBrandIsExist(value any) error {
	var id string
	switch v := value.(type) {
	case string:
		id = v
	case *string:
		if v == nil {
			return nil
		}
		id = *v
	default:
		return richerror.New("brandvalidator.mediaBrandIsExist.type_assertion").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("global.validation_failed"))
	}

	if id == "" {
		return nil
	}

	_, err := v.mediaSvc.FindMediaByID(context.Background(), mediaparam.GetMediaByIDRequest{ID: id})

	if err != nil {
		var rErr *richerror.RichError
		if errors.As(err, &rErr) {
			if rErr.Kind() == richerror.KindNotFound {
				return errors.New(rErr.MessageKey())
			}
		}
		return err
	}
	return nil
}

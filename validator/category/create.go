package categoryvalidator

import (
	"context"
	"errors"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateCreateCategoryRequest(req categoryparam.CreateCategoryRequest) error {
	const op = "backofficevalidatorauth.ValidateLoginWithPhoneNumberRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.Required.Error("validation.required")),
		validation.Field(&req.ParentID,
			validation.When(req.ParentID != "", validation.Required.Error("validation.required")),
			validation.By(v.categoryIsExist),
		),
		validation.Field(&req.MediaID,
			validation.Required.Error("validation.required"),
			validation.By(v.mediaIsExist),
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

func (v Validator) mediaIsExist(value any) error {
	mediaID, ok := value.(string)
	if !ok {
		return richerror.New("categoryvalidator.mediaIsExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}

	_, err := v.mediaSvc.FindMediaByID(context.Background(), mediaparam.GetMediaByIDRequest{ID: mediaID})

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

func (v Validator) categoryIsExist(value any) error {
	id, ok := value.(string)
	if !ok {
		return richerror.New("categoryvalidator.categoryIsExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}

	category, err := v.repo.FindByID(context.Background(), id)

	if err != nil || category.ID == "" {
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

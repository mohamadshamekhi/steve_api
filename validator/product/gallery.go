package productvalidator

import (
	"context"
	"errors"

	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateAddProductGalleryRequest(req productparam.AddProductGalleryRequest) error {
	const op = "productvalidator.ValidateAddProductGalleryRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.MediaID,
			validation.Required.Error("validation.required"), validation.By(v.mediaIsExist)),
		validation.Field(&req.MediaURL,
			validation.Required.Error("validation.required")),
		validation.Field(&req.ProductID, validation.Required.Error("validation.required"), validation.By(v.checkProductIsExist)),
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

func (v Validator) ValidateDeleteProductGalleryRequest(req productparam.DeleteProductGalleryRequest) error {
	const op = "productvalidator.ValidateDeleteProductGalleryRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.MediaID,
			validation.Required.Error("validation.required"), validation.By(v.mediaIsExist)),
		validation.Field(&req.ProductID, validation.Required.Error("validation.required"), validation.By(v.checkProductIsExist)),
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
func (v Validator) ValidateSetPrimaryGalleryRequest(req productparam.SetPrimaryProductGalleryRequest) error {
	const op = "productvalidator.ValidateSetPrimaryGalleryRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.MediaID,
			validation.Required.Error("validation.required"), validation.By(v.mediaIsExist)),
		validation.Field(&req.ProductID, validation.Required.Error("validation.required"), validation.By(v.checkProductIsExist)),
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

func (v Validator) mediaIsExist(value any) error {
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
		return richerror.New("brandvalidator.BrandIsUnique.type_assertion").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type for brand name"))
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

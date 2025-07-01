package mediavalidator

import (
	"errors"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUploadRequest(req mediaparam.UploadMediaRequest) error {
	const op = "mediavalidator.ValidateUploadRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Domain, validation.Required),
		validation.Field(&req.UploaderID, validation.Required),
		validation.Field(&req.UploaderType, validation.In("backoffice", "user")),
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

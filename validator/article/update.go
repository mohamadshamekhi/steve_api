package articlevalidator

import (
	"errors"
	articleparam "github.com/abolfazlnorzad/steve/param/article"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateArticleRequest(req articleparam.UpdateArticleRequest) error {
	const op = "categoryvalidator.ValidateUpdateCategoryRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Title,
			validation.When(req.Title != nil,
				validation.Required.Error("validation.required")),
		),
		validation.Field(&req.Body,
			validation.When(req.Body != nil,
				validation.Required.Error("validation.required"),
			),
		),
		validation.Field(&req.MediaID,
			validation.When(req.MediaID != nil,
				validation.Required.Error("validation.required"),
				validation.By(v.mediaIsExist),
			),
		),
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

package articlevalidator

import (
	"errors"

	articleparam "github.com/abolfazlnorzad/steve/param/article"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateArticleRequest(req articleparam.UpdateArticleRequest) error {
	const op = "articlevalidator.ValidateUpdateArticleRequest"

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
		validation.Field(&req.Slug,
			validation.When(req.Slug != nil,
				validation.Required.Error("validation.required"),
				validation.By(v.SlugIsUnique),
			),
		),
		validation.Field(&req.MediaID,
			validation.When(req.MediaID != nil,
				validation.Required.Error("validation.required"),
				validation.By(v.mediaIsExist),
			),
		),
		validation.Field(&req.Tags,
			validation.When(req.Tags != nil,
				validation.Required.Error("validation.required"),
			),
		),
		validation.Field(&req.MetaDescription,
			validation.When(req.MetaDescription != nil,
				validation.Required.Error("validation.required"),
			),
		),
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

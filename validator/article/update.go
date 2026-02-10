package articlevalidator

import (
	"context"
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
				validation.By(v.SlugIsUniqueForUpdate(req.ID)),
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

func (v Validator) SlugIsUniqueForUpdate(id string) validation.RuleFunc {
	return func(value any) error {
		var slug string
		switch v := value.(type) {
		case string:
			slug = v
		case *string:
			if v == nil {
				return nil
			}
			slug = *v
		default:
			return richerror.New("articlevalidator.SlugIsUnique.type_assertion").
				WithKind(richerror.KindUnexpected).
				WithWarpError(errors.New("unexpected type for slug"))
		}

		if slug == "" {
			return nil
		}

		article, err := v.aRepo.FindBySlug(context.Background(), slug)
		if err != nil {
			var rErr *richerror.RichError
			if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
				return nil
			}
			return err
		}

		if article.ID == id {
			return nil
		}

		return errors.New("validation.slug_already_exists")
	}
}

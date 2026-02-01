package articlevalidator

import (
	"context"
	"errors"

	articleparam "github.com/abolfazlnorzad/steve/param/article"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateCreateArticleRequest(req articleparam.AddArticleRequest) error {
	const op = "articlevalidator.ValidateCreateArticleRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Title,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Body,
			validation.Required.Error("validation.required"),
		),
		validation.Field(&req.Slug,
			validation.Required.Error("validation.required"),
			validation.By(v.SlugIsUnique),
		),
		validation.Field(&req.MediaID,
			validation.Required.Error("validation.required"),
			validation.By(v.mediaIsExist),
		)); err != nil {
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

func (v Validator) SlugIsUnique(value any) error {
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

	_, err := v.aRepo.FindBySlug(context.Background(), slug)
	if err != nil {
		var rErr *richerror.RichError
		if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
			return nil
		}
		return err
	}

	return errors.New("validation.slug_already_exists")
}

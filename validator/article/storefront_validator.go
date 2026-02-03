package articlevalidator

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	articleparam "github.com/abolfazlnorzad/steve/param/article"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Reader interface {
	FindByID(ctx context.Context, id string) (entity.Article, error)
	FindBySlug(ctx context.Context, slug string) (entity.Article, error)
}

type StorefrontValidator struct {
	repo Reader
}

func NewStorefrontValidator(repo Reader) StorefrontValidator {
	return StorefrontValidator{repo: repo}
}

func (v StorefrontValidator) isArticleExist(value any) error {
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
		return richerror.New("articlevalidator.articleIsUnique.type_assertion").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type for article id"))
	}

	if id == "" {
		return errors.New("id is empty")
	}

	a, err := v.repo.FindByID(context.Background(), id)

	if err != nil || a.ID == "" {
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

func (v StorefrontValidator) isArticleExistBySlug(value any) error {
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
		return richerror.New("articlevalidator.articleIsUnique.type_assertion").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type for article slug"))
	}

	if slug == "" {
		return errors.New("slug is empty")
	}
	a, err := v.repo.FindBySlug(context.Background(), slug)

	if err != nil || a.ID == "" {
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

func (v StorefrontValidator) ValidateGetStoreArticleRequest(req articleparam.GetStorefrontArticleRequest) error {
	const op = "articlevalidator.ValidateGetStoreArticleRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID,
			validation.Required.Error("validation.required"),
			validation.By(v.isArticleExist)),
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

func (v StorefrontValidator) ValidateGetBySlugStoreArticleRequest(req articleparam.GetBySlugStorefrontArticleRequest) error {
	const op = "articlevalidator.ValidateGetBySlugStoreArticleRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Slug,
			validation.By(func(value interface{}) error {
				if value == "" || value == ":slug" {
					return errors.New("slug is empty")
				}
				return nil
			}),

			validation.By(v.isArticleExistBySlug)),
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
	return nil
}

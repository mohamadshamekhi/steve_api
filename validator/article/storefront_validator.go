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
		return richerror.New("brandvalidator.BrandIsUnique.type_assertion").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type for brand name"))
	}

	if id == "" {
		return nil
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

func (v StorefrontValidator) ValidateGetStoreArticleRequest(req articleparam.GetStorefrontArticleRequest) error {
	const op = "backofficevalidatorauth.ValidateGetStoreArticleRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID,
			validation.Required.Error("validation.required"),
			validation.By(v.isArticleExist)),
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

package commentvalidator

import (
	"context"
	"errors"

	commentparam "github.com/abolfazlnorzad/steve/param/comment"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateAddCommentRequest(req commentparam.AddCommentRequest) error {
	const op = "commentvalidator.ValidateAddCommentRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Body, validation.Required,
			validation.Required.Error("validation.required")),
		validation.Field(&req.CommentableType, validation.Required,
			validation.Required.Error("validation.required"), validation.In("article", "product")),
		validation.Field(&req.CommentableID, validation.Required,
			validation.Required.Error("validation.required")),
		validation.Field(&req.ParentID,
			validation.When(req.ParentID != nil, validation.Required.Error("validation.required"), validation.By(v.commentIsExist)),
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

func (v Validator) commentIsExist(value any) error {
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

	res, err := v.repo.CommentIsExist(context.Background(), id)
	if err != nil {
		return err
	}
	if res {
		return nil
	}
	return errors.New("global.not-found")
}

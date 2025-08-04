package commentvalidator

import (
	"errors"
	commentparam "github.com/abolfazlnorzad/steve/param/comment"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUpdateCommentRequest(req commentparam.UpdateCommentRequest) error {
	const op = "commentvalidator.ValidateAddCommentRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.CommentID,
			validation.Required.Error("validation.required"), validation.By(v.commentIsExist))); err != nil {

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

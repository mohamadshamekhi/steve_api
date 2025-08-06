package mediavalidator

import (
	"context"
	"errors"
	"fmt"

	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateGetFile(req mediaparam.GetUrlRequest) error {
	const op = "mediavalidator.ValidateGetFile"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Media.ID, validation.Required, validation.By(v.checkMediaIsExistByID)),
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

func (v Validator) checkMediaIsExistByID(value any) error {
	mediaID, ok := value.(string)
	if !ok {
		return fmt.Errorf(errmsg.ErrorMsgSomethingWentWrong)
	}

	res, err := v.repo.CheckMediaIsExist(context.Background(), mediaID)
	if err != nil {
		return err
	}
	if !res {
		return errors.New(errmsg.ErrorMsgNotFound)
	}

	return nil
}

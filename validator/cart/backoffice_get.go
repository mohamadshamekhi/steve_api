package cartvalidator

import (
	"context"
	"errors"

	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateGetCartRequest(req cartparam.GetCartRequest) error {
	const op = "cartvalidator.ValidateGetCartRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ID,
			validation.Required.Error("validation.required"),
			validation.By(v.cartExist)),
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

func (v Validator) cartExist(value any) error {
	id, ok := value.(string)
	if !ok {
		return richerror.New("cartvalidator.cartExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	_, err := v.repo.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}

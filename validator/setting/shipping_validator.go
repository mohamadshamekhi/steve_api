package settingvalidator

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	settingparam "github.com/abolfazlnorzad/steve/param/setting"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Repo interface {
	FindByID(ctx context.Context, id string) (entity.Shipping, error)
}

type ShippingValidator struct {
	repo Repo
}

func NewShippingValidator(repo Repo) ShippingValidator {
	return ShippingValidator{repo: repo}
}

func (v ShippingValidator) ValidateAddShippingRequest(req settingparam.AddShippingRequest) error {
	const op = "settingvalidator.ValidateAddShippingRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.Required.Error("validation.required")),
		validation.Field(&req.IsActive,
			validation.NotNil.Error("validation.required")),
		validation.Field(&req.Cost,
			validation.Required.Error("validation.required"),
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

func (v ShippingValidator) ValidateUpdateShippingRequest(req settingparam.UpdateShippingRequest) error {
	const op = "settingvalidator.ValidateUpdateShippingRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.When(req.Name != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Cost,
			validation.When(req.Cost != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.ShippingID,
			validation.Required.Error("validation.required"), validation.By(v.checkShippingIsExist)),
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

func (v ShippingValidator) ValidateFindByIDRequest(req settingparam.FindByIDShippingRequest) error {
	const op = "settingvalidator.ValidateFindByIDRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ShippingID,
			validation.Required.Error("validation.required"), validation.By(v.checkShippingIsExist)),
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

func (v ShippingValidator) checkShippingIsExist(value any) error {
	const op = "shippingvalidator.checkShippingIsExist"
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
		return richerror.New("brandvalidator.mediaBrandIsExist.type_assertion").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("global.validation_failed"))
	}
	_, err := v.repo.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}

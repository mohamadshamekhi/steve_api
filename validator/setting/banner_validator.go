package settingvalidator

import (
	"errors"

	settingparam "github.com/abolfazlnorzad/steve/param/setting"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type BannerValidator struct {
}

func NewBannerValidator() BannerValidator {
	return BannerValidator{}
}

func (v BannerValidator) ValidateAddBannerRequest(req settingparam.CreateBannerRequest) error {
	const op = "settingvalidator.ValidateAddBannerRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Title,
			validation.Required.Error("validation.required")),
		validation.Field(&req.IsActive,
			validation.NotNil.Error("validation.required")),
		validation.Field(&req.MediaID,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Link,
			validation.Required.Error("validation.required")),
		validation.Field(&req.Location,
			validation.Required.Error("validation.required")),
		validation.Field(&req.AltText,
			validation.Required.Error("validation.required")),
		validation.Field(&req.SortOrder,
			validation.Required.Error("validation.required")),

		validation.Field(&req.ExpiresAt,
			validation.Min(req.StartsAt).Error("تاریخ انقضا باید بعد از تاریخ شروع باشد"),
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

func (v BannerValidator) ValidateUpdateBannerRequest(req settingparam.UpdateBannerRequest) error {
	const op = "settingvalidator.ValidateUpdateBannerRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Title,
			validation.When(req.Title != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.AltText,
			validation.When(req.AltText != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.Link,
			validation.When(req.Link != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.MediaID,
			validation.When(req.MediaID != nil, validation.Required.Error("validation.required"))),
		validation.Field(&req.SortOrder,
			validation.When(req.SortOrder != nil, validation.Required.Error("validation.required"))),
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

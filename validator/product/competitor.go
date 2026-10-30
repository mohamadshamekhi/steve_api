package productvalidator

import (
	"errors"
	"regexp"
	"time"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// TorobURLRegex accepts only a Torob product page: https://torob.com/p/<uuid>
// with an optional trailing slash, nothing after it.
var TorobURLRegex = regexp.MustCompile(`^https://torob\.com/p/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/?$`)

func (v Validator) ValidateUpdateCompetitorPriceRequest(req productparam.UpdateCompetitorPriceRequest) error {
	const op = "productvalidator.ValidateUpdateCompetitorPriceRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ProductID, validation.Required.Error("validation.required")),
		validation.Field(&req.MaxPrice,
			validation.Min(req.MinPrice).Error("validation.invalid_format")),
		validation.Field(&req.CheckedAt,
			validation.NotNil.Error("validation.required"),
			validation.By(notInFuture)),
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

// notInFuture allows a few minutes of clock skew between the agent and us.
func notInFuture(value any) error {
	t, ok := value.(*time.Time)
	if !ok || t == nil {
		return nil
	}
	if t.After(time.Now().Add(5 * time.Minute)) {
		return errors.New("validation.invalid_format")
	}
	return nil
}

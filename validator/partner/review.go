package partnervalidator

import (
	partnerparam "github.com/abolfazlnorzad/steve/param/partner"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateApproveRequest(req partnerparam.ApproveRequest) error {
	const op = "partnervalidator.ValidateApproveRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.DiscountPercent,
			validation.Max(uint(90)).Error("partner.invalid_discount_percent")),
	); err != nil {
		return wrap(op, err)
	}
	return nil
}

func (v Validator) ValidateRejectRequest(req partnerparam.RejectRequest) error {
	const op = "partnervalidator.ValidateRejectRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Reason,
			validation.Required.Error("validation.required"),
			validation.Length(3, 500).Error("validation.invalid_format")),
	); err != nil {
		return wrap(op, err)
	}
	return nil
}

func (v Validator) ValidateUpdatePricingRequest(req partnerparam.UpdatePricingRequest) error {
	const op = "partnervalidator.ValidateUpdatePricingRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.DiscountPercent,
			validation.Max(uint(90)).Error("partner.invalid_discount_percent")),
	); err != nil {
		return wrap(op, err)
	}
	return nil
}

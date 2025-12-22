package discountparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
)

type UpdateDiscountRequest struct {
	ID                    string                        `json:"id"`
	Value                 *uint                         `json:"value"`
	MinimumPurchaseAmount *uint                         `json:"minimum_purchase_amount"`
	ExpiresAt             *time.Time                    `json:"expires_at"`
	TotalUsageLimit       *int                          `json:"total_usage_limit"`
	UsageLimitPerUser     *int                          `json:"usage_limit_per_user"`
	IsSpecial             *bool                         `json:"is_special"`
	ApplicableTo          *string                       `json:"applicable_to"`
	ApplicableProductIDs  *[]entity.SpecialDiscountItem `json:"applicable_product_ids,omitempty"`
	ApplicableCategoryIDs *[]string                     `json:"applicable_category_ids"`
	ApplicableUserIDs     *[]string                     `json:"applicable_user_ids"`
	IsActive              *bool                         `json:"is_active"`
}

type UpdateDiscountResponse struct {
}

type UpdateDiscountUsageAfterRegisterOrderRequest struct {
	AppliedDiscount []entity.AppliedDiscount
}

type UpdateDiscountUsageAfterRegisterOrderResponse struct {
}

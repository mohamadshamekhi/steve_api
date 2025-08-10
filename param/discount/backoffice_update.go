package discountparam

import "time"

type UpdateDiscountRequest struct {
	ID                    string
	Value                 *uint      `json:"value"`
	MinimumPurchaseAmount *uint      `json:"minimum_purchase_amount"`
	ExpiresAt             *time.Time `json:"expires_at"`
	TotalUsageLimit       *int       `json:"total_usage_limit"`
	UsageLimitPerUser     *int       `json:"usage_limit_per_user"`
	ApplicableProductIDs  *[]string  `json:"applicable_product_ids"`
	ApplicableCategoryIDs *[]string  `json:"applicable_category_ids"`
	ApplicableUserIDs     *[]string  `json:"applicable_user_ids"`
	IsActive              *bool      `json:"is_active"`
}

type UpdateDiscountResponse struct {
}

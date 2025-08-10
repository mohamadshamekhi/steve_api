package discountparam

import "time"

type AddDiscountRequest struct {
	Code                  string    `json:"code,omitempty"`
	Type                  string    `bson:"type"`
	Value                 uint      `json:"value"`
	MinimumPurchaseAmount uint      `json:"minimum_purchase_amount,omitempty"`
	StartsAt              time.Time `json:"starts_at"`
	ExpiresAt             time.Time `json:"expires_at"`
	TotalUsageLimit       int       `json:"total_usage_limit,omitempty"`
	ApplicableProductIDs  []string  `json:"applicable_product_ids,omitempty"`
	ApplicableCategoryIDs []string  `json:"applicable_category_ids,omitempty"`
	ApplicableUserIDs     []string  `json:"applicable_user_ids,omitempty"`
	IsActive              bool      `json:"is_active"`
}

type AddDiscountResponse struct {
}

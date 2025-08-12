package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type DiscountType string

const (
	PercentageDiscount  DiscountType = "percentage"
	FixedAmountDiscount DiscountType = "fixed_amount"
)

type Discount struct {
	ID   string       `bson:"_id" json:"id"`
	Code string       `bson:"code,omitempty" json:"code"`
	Type DiscountType `bson:"type" json:"type"`

	Value uint `bson:"value" json:"value"`

	MinimumPurchaseAmount uint `bson:"minimum_purchase_amount,omitempty"json:"minimum_purchase_amount,omitempty"`

	StartsAt time.Time `bson:"starts_at" json:"starts_at ,omitempty"`

	ExpiresAt time.Time `bson:"expires_at" json:"expires_at ,omitempty"`

	TotalUsageLimit int `bson:"total_usage_limit,omitempty" json:"total_usage_limit,omitempty"`

	UsageCount int `bson:"usage_count" json:"usage_count,omitempty"`

	UsageLimitPerUser int `bson:"usage_limit_per_user,omitempty" json:"usage_limit_per_user,omitempty"`

	ApplicableProductIDs []string `bson:"applicable_product_ids,omitempty" json:"applicable_product_ids,omitempty"`

	ApplicableCategoryIDs []string `bson:"applicable_category_ids,omitempty" json:"applicable_category_ids,omitempty"`

	ApplicableUserIDs []string `bson:"applicable_user_ids,omitempty" json:"applicable_user_ids,omitempty"`

	IsActive  bool      `bson:"is_active" json:"is_active"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (d *Discount) BeforeCreate() {
	d.ID = ulid.Make().String()
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	d.UsageCount = 0
}

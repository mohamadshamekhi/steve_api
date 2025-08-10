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
	ID   string       `bson:"_id"`
	Code string       `bson:"code,omitempty"`
	Type DiscountType `bson:"type"`

	Value uint `bson:"value"`

	MinimumPurchaseAmount uint `bson:"minimum_purchase_amount,omitempty"`

	StartsAt time.Time `bson:"starts_at"`

	ExpiresAt time.Time `bson:"expires_at"`

	TotalUsageLimit int `bson:"total_usage_limit,omitempty"`

	UsageCount int `bson:"usage_count"`

	UsageLimitPerUser int `bson:"usage_limit_per_user,omitempty"`

	ApplicableProductIDs []string `bson:"applicable_product_ids,omitempty"`

	ApplicableCategoryIDs []string `bson:"applicable_category_ids,omitempty"`

	ApplicableUserIDs []string `bson:"applicable_user_ids,omitempty"`

	IsActive  bool      `bson:"is_active"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (d *Discount) BeforeCreate() {
	d.ID = ulid.Make().String()
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	d.IsActive = true
	d.UsageCount = 0
}

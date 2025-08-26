package entity

import (
	"time"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
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

func (d *Discount) CheckIsValidForUser(userID string, amount uint) error {
	const op = "discount.CheckIsValidForUser"

	if !d.IsActive {
		return richerror.New(op).WithKind(richerror.KindBadRequest).
			WithMessageKey("discount.is_not_active")
	}

	if !d.StartsAt.IsZero() && d.StartsAt.After(time.Now()) {
		return richerror.New(op).WithKind(richerror.KindBadRequest).
			WithMessageKey("discount.not_yet_started")
	}

	if !d.ExpiresAt.IsZero() && d.ExpiresAt.Before(time.Now()) {
		return richerror.New(op).WithKind(richerror.KindBadRequest).
			WithMessageKey("discount.has_expired")
	}

	if d.TotalUsageLimit > 0 && d.UsageCount >= d.TotalUsageLimit {
		return richerror.New(op).WithKind(richerror.KindBadRequest).
			WithMessageKey("discount.usage_limit_reached")
	}

	if d.MinimumPurchaseAmount > 0 && amount < d.MinimumPurchaseAmount {
		return richerror.New(op).WithKind(richerror.KindBadRequest).
			WithMessageKey("discount.minimum_purchase_not_met")
	}

	if len(d.ApplicableUserIDs) > 0 {
		userIsAllowed := false
		for _, id := range d.ApplicableUserIDs {
			if userID == id {
				userIsAllowed = true
				break
			}
		}
		if !userIsAllowed {
			return richerror.New(op).WithKind(richerror.KindBadRequest).
				WithMessageKey("discount.not_applicable_for_user")
		}
	}

	return nil
}

package entity

import "time"

// PriceReviewSetting is a singleton document holding the admin-configured
// interval (in days) after which a product's price must be re-checked.
type PriceReviewSetting struct {
	ID           string    `bson:"_id"`
	IntervalDays int       `bson:"interval_days"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

// PriceReviewSettingID is the fixed document id for the singleton setting.
const PriceReviewSettingID = "default"

// DefaultPriceReviewIntervalDays is used when no setting document exists yet.
const DefaultPriceReviewIntervalDays = 20

package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// BirthdaySettingsID is the fixed _id of the single birthday-settings document.
const BirthdaySettingsID = "singleton"

// BirthdaySettings is the admin-editable configuration of the birthday
// discount program. One document (_id == BirthdaySettingsID); when it doesn't
// exist the service falls back to config defaults.
type BirthdaySettings struct {
	ID string `bson:"_id" json:"-"`

	Enabled bool `bson:"enabled" json:"enabled"`
	// DaysBefore is how many days ahead of the birthday the code is sent
	// (0 = on the birthday itself).
	DaysBefore int `bson:"days_before" json:"days_before"`
	// SendHour is the Tehran hour (0-23) from which the day's codes and SMS go out.
	SendHour              int    `bson:"send_hour" json:"send_hour"`
	DiscountType          string `bson:"discount_type" json:"discount_type"` // "percentage" | "fixed_amount"
	DiscountValue         uint   `bson:"discount_value" json:"discount_value"`
	ValidDays             int    `bson:"valid_days" json:"valid_days"`
	MinimumPurchaseAmount uint   `bson:"minimum_purchase_amount" json:"minimum_purchase_amount"`

	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	UpdatedBy string    `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
}

// BirthdayReward records the discount code granted to a user for one
// birthday. Unique on (user_id, year) so a user gets at most one per year.
type BirthdayReward struct {
	ID           string    `bson:"_id" json:"id"`
	UserID       string    `bson:"user_id" json:"user_id"`
	UserName     string    `bson:"user_name,omitempty" json:"user_name,omitempty"`
	PhoneNumber  string    `bson:"phone_number" json:"phone_number"`
	Year         int       `bson:"year" json:"year"`
	BirthdayOn   time.Time `bson:"birthday_on" json:"birthday_on"`
	DiscountID   string    `bson:"discount_id,omitempty" json:"discount_id,omitempty"`
	DiscountCode string    `bson:"discount_code,omitempty" json:"discount_code,omitempty"`
	SMSSent      bool      `bson:"sms_sent" json:"sms_sent"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

func (r *BirthdayReward) BeforeCreate() {
	r.ID = ulid.Make().String()
	r.CreatedAt = time.Now()
}

// MonthDay is a calendar day independent of year, used to match birth dates.
type MonthDay struct {
	Month time.Month
	Day   int
}

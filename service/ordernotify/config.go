package ordernotifyservice

import "time"

type Config struct {
	// Enabled turns every order SMS (thank-you and status change) on or off.
	Enabled bool `koanf:"enabled"`

	// Post-purchase thank-you coupon, scoped to the buyer and single-use.
	ThankYouDiscountEnabled       bool   `koanf:"thank_you_discount_enabled"`
	ThankYouDiscountType          string `koanf:"thank_you_discount_type"` // percentage | fixed_amount
	ThankYouDiscountValue         uint   `koanf:"thank_you_discount_value"`
	ThankYouValidDays             int    `koanf:"thank_you_valid_days"`
	ThankYouMinimumPurchaseAmount uint   `koanf:"thank_you_minimum_purchase_amount"`

	SMSTimeout time.Duration `koanf:"sms_timeout"`
}

func (c Config) smsTimeout() time.Duration {
	if c.SMSTimeout <= 0 {
		return 30 * time.Second
	}
	return c.SMSTimeout
}

func (c Config) validDays() int {
	if c.ThankYouValidDays <= 0 {
		return 30
	}
	return c.ThankYouValidDays
}

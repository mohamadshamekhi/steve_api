package referralservice

import "github.com/abolfazlnorzad/steve/entity"

// Config drives the "invite friends" program. Loaded via koanf with the
// ST_REFERRAL__ env prefix — spell out every nested level, e.g.
// ST_REFERRAL__REWARD_THRESHOLD=5, ST_REFERRAL__INVITE_BASE_URL=https://...
type Config struct {
	// Enabled turns the whole referral program on/off.
	Enabled bool `koanf:"enabled"`
	// InviteBaseURL is the storefront page that reads the ?ref=<code> query
	// param; the shareable link is InviteBaseURL + "?ref=" + code.
	InviteBaseURL string `koanf:"invite_base_url"`
	// RewardThreshold is how many joins earn the referrer one discount code.
	RewardThreshold int `koanf:"reward_threshold"`

	// --- reward code granted to the referrer every RewardThreshold joins ---
	RewardDiscountType  string `koanf:"reward_discount_type"` // "percentage" | "fixed_amount"
	RewardDiscountValue uint   `koanf:"reward_discount_value"`
	RewardValidDays     int    `koanf:"reward_valid_days"`
	// RewardMinimumPurchaseAmount optionally gates the reward code.
	RewardMinimumPurchaseAmount uint `koanf:"reward_minimum_purchase_amount"`

	// --- optional welcome code granted to the invitee on join ---
	WelcomeDiscountEnabled bool   `koanf:"welcome_discount_enabled"`
	WelcomeDiscountType    string `koanf:"welcome_discount_type"`
	WelcomeDiscountValue   uint   `koanf:"welcome_discount_value"`
	WelcomeValidDays       int    `koanf:"welcome_valid_days"`
}

func (c Config) rewardThreshold() int {
	if c.RewardThreshold <= 0 {
		return 5
	}
	return c.RewardThreshold
}

func (c Config) rewardValidDays() int {
	if c.RewardValidDays <= 0 {
		return 30
	}
	return c.RewardValidDays
}

func (c Config) welcomeValidDays() int {
	if c.WelcomeValidDays <= 0 {
		return 30
	}
	return c.WelcomeValidDays
}

// applySettings overlays an admin-saved settings document onto the config
// defaults. Enabled and InviteBaseURL stay deployment-controlled (config/env
// only); everything about the reward/welcome coupons becomes admin-editable.
func (c Config) applySettings(s entity.ReferralSettings) Config {
	c.RewardThreshold = s.RewardThreshold
	c.RewardDiscountType = s.RewardDiscountType
	c.RewardDiscountValue = s.RewardDiscountValue
	c.RewardValidDays = s.RewardValidDays
	c.RewardMinimumPurchaseAmount = s.RewardMinimumPurchaseAmount
	c.WelcomeDiscountEnabled = s.WelcomeDiscountEnabled
	c.WelcomeDiscountType = s.WelcomeDiscountType
	c.WelcomeDiscountValue = s.WelcomeDiscountValue
	c.WelcomeValidDays = s.WelcomeValidDays
	return c
}

// toSettings projects the (effective) config into a settings document — used to
// seed the admin editor with current values before any customisation.
func (c Config) toSettings() entity.ReferralSettings {
	return entity.ReferralSettings{
		RewardThreshold:             c.rewardThreshold(),
		RewardDiscountType:          c.RewardDiscountType,
		RewardDiscountValue:         c.RewardDiscountValue,
		RewardValidDays:             c.rewardValidDays(),
		RewardMinimumPurchaseAmount: c.RewardMinimumPurchaseAmount,
		WelcomeDiscountEnabled:      c.WelcomeDiscountEnabled,
		WelcomeDiscountType:         c.WelcomeDiscountType,
		WelcomeDiscountValue:        c.WelcomeDiscountValue,
		WelcomeValidDays:            c.welcomeValidDays(),
	}
}

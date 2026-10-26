package birthdayservice

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
)

// Config holds the defaults of the birthday-discount program, used until an
// admin saves settings in the backoffice. Loaded via koanf with the
// ST_BIRTHDAY__ env prefix, e.g. ST_BIRTHDAY__DAYS_BEFORE=3.
type Config struct {
	Enabled               bool   `koanf:"enabled"`
	DaysBefore            int    `koanf:"days_before"`
	SendHour              int    `koanf:"send_hour"`     // 0-23, Tehran time
	DiscountType          string `koanf:"discount_type"` // "percentage" | "fixed_amount"
	DiscountValue         uint   `koanf:"discount_value"`
	ValidDays             int    `koanf:"valid_days"`
	MinimumPurchaseAmount uint   `koanf:"minimum_purchase_amount"`
}

func (c Config) toSettings() entity.BirthdaySettings {
	return entity.BirthdaySettings{
		Enabled:               c.Enabled,
		DaysBefore:            c.DaysBefore,
		SendHour:              c.SendHour,
		DiscountType:          c.DiscountType,
		DiscountValue:         c.DiscountValue,
		ValidDays:             c.ValidDays,
		MinimumPurchaseAmount: c.MinimumPurchaseAmount,
	}
}

const maxDaysBefore = 30

func daysBefore(s entity.BirthdaySettings) int {
	switch {
	case s.DaysBefore < 0:
		return 0
	case s.DaysBefore > maxDaysBefore:
		return maxDaysBefore
	}
	return s.DaysBefore
}

func sendHour(s entity.BirthdaySettings) int {
	if s.SendHour < 0 || s.SendHour > 23 {
		return 10
	}
	return s.SendHour
}

// rewardCooldown is the minimum gap between two birthday codes for one user,
// whatever their birth date says — guards against birth-date edits that move
// the birthday into another year.
const rewardCooldown = 330 * 24 * time.Hour

func validDays(s entity.BirthdaySettings) int {
	if s.ValidDays <= 0 {
		return 10
	}
	return s.ValidDays
}

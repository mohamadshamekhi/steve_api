package birthdayparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

// SettingsBody is the editable birthday-discount configuration, shared by the
// GET response and the PUT request (PUT fully replaces it).
type SettingsBody struct {
	Enabled               bool   `json:"enabled"`
	DaysBefore            int    `json:"days_before"`
	SendHour              int    `json:"send_hour"`     // 0-23, Tehran time
	DiscountType          string `json:"discount_type"` // "percentage" | "fixed_amount"
	DiscountValue         uint   `json:"discount_value"`
	ValidDays             int    `json:"valid_days"`
	MinimumPurchaseAmount uint   `json:"minimum_purchase_amount"`
}

type GetSettingsResponse struct {
	Settings SettingsBody `json:"settings"`
	// Source is "config" while no admin override exists yet, then "custom".
	Source    string     `json:"source"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	UpdatedBy string     `json:"updated_by,omitempty"`
}

// UpdateSettingsRequest is the PUT body — the settings fields flat.
type UpdateSettingsRequest struct {
	AdminUserID string `json:"-"`
	SettingsBody
}

type UpdateSettingsResponse struct {
	Settings SettingsBody `json:"settings"`
}

type ListRewardsRequest struct {
	Pagination commonparam.PaginationRequest
}

type ListRewardsResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Items      []entity.BirthdayReward        `json:"items"`
}

// RunResponse summarises one pass of the birthday job.
type RunResponse struct {
	Enabled bool `json:"enabled"`
	Matched int  `json:"matched"`
	Granted int  `json:"granted"`
	Skipped int  `json:"skipped"`
	Failed  int  `json:"failed"`
}

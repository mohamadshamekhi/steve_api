package referralparam

import (
	"time"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

// AdminListReferralsRequest is the backoffice referral leaderboard query.
type AdminListReferralsRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" query:"pagination"`
	// MinJoined filters to accounts with at least this many completed joins.
	MinJoined int `json:"min_joined" query:"min_joined"`
}

// AdminAccountRow is one referrer in the admin list / the header of the detail view.
type AdminAccountRow struct {
	UserID          string    `json:"user_id"`
	ReferrerName    string    `json:"referrer_name,omitempty"`
	ReferrerPhone   string    `json:"referrer_phone,omitempty"`
	Code            string    `json:"code"`
	TotalJoined     int       `json:"total_joined"`
	RewardedCount   int       `json:"rewarded_count"`
	PendingProgress int       `json:"pending_progress"`
	RewardsEarned   int       `json:"rewards_earned"`
	CreatedAt       time.Time `json:"created_at"`
}

type AdminListReferralsResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Items      []AdminAccountRow              `json:"items"`
}

type AdminGetReferralRequest struct {
	UserID     string                        `json:"user_id"`
	Pagination commonparam.PaginationRequest `json:"pagination" query:"pagination"`
}

// AdminGetReferralResponse is the per-referrer detail: their account header, a
// paginated list of their invitees, and every reward code they earned.
type AdminGetReferralResponse struct {
	Account    AdminAccountRow                `json:"account"`
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Referrals  []ReferralItem                 `json:"referrals"`
	Rewards    []RewardItem                   `json:"rewards"`
}

// ReferralSettingsBody is the editable reward/welcome-coupon configuration,
// shared by the GET response and the PUT request (PUT fully replaces it).
type ReferralSettingsBody struct {
	RewardThreshold             int    `json:"reward_threshold"`
	RewardDiscountType          string `json:"reward_discount_type"` // "percentage" | "fixed_amount"
	RewardDiscountValue         uint   `json:"reward_discount_value"`
	RewardValidDays             int    `json:"reward_valid_days"`
	RewardMinimumPurchaseAmount uint   `json:"reward_minimum_purchase_amount"`

	WelcomeDiscountEnabled bool   `json:"welcome_discount_enabled"`
	WelcomeDiscountType    string `json:"welcome_discount_type"`
	WelcomeDiscountValue   uint   `json:"welcome_discount_value"`
	WelcomeValidDays       int    `json:"welcome_valid_days"`
}

type GetReferralSettingsResponse struct {
	Settings ReferralSettingsBody `json:"settings"`
	// Source is "config" while no admin override exists yet, then "custom".
	Source    string     `json:"source"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	UpdatedBy string     `json:"updated_by,omitempty"`
	// ProgramEnabled / InviteBaseURL are deployment-controlled (read-only here).
	ProgramEnabled bool   `json:"program_enabled"`
	InviteBaseURL  string `json:"invite_base_url"`
}

// UpdateReferralSettingsRequest is the PUT body — the settings fields flat, no
// wrapper object.
type UpdateReferralSettingsRequest struct {
	AdminUserID string `json:"-"`
	ReferralSettingsBody
}

type UpdateReferralSettingsResponse struct {
	Settings ReferralSettingsBody `json:"settings"`
}

// AdminReferralStatsResponse is the program-wide dashboard summary.
type AdminReferralStatsResponse struct {
	TotalAccounts  int64   `json:"total_accounts"`
	TotalInvites   int64   `json:"total_invites"`
	TotalJoined    int64   `json:"total_joined"`
	TotalRewards   int64   `json:"total_rewards"`
	ConversionRate float64 `json:"conversion_rate"`
}

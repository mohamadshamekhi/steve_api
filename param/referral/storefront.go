package referralparam

import (
	"time"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

// GetMyLinkResponse powers the "invite friends" screen.
type GetMyLinkResponse struct {
	Code string `json:"code"`
	// URL is the shareable link (InviteBaseURL + ?ref=<code>).
	URL string `json:"url"`
	// TotalJoined is how many people registered through this user's link.
	TotalJoined int `json:"total_joined"`
	// PendingProgress is joined invitees not yet turned into a reward.
	PendingProgress int `json:"pending_progress"`
	// RewardThreshold is how many joins earn one discount code.
	RewardThreshold int `json:"reward_threshold"`
	// NextRewardIn is how many more joins until the next discount code.
	NextRewardIn int `json:"next_reward_in"`
	// RewardsEarned is the total number of discount codes granted so far.
	RewardsEarned int `json:"rewards_earned"`
}

// InviteByPhoneRequest asks us to text an invite to a phone number.
type InviteByPhoneRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type InviteByPhoneResponse struct {
	ReferralID string `json:"referral_id"`
	// URL is the shareable invite link — always returned so the frontend can
	// offer a "copy link" fallback when the SMS could not be delivered.
	URL string `json:"url"`
	// SMSSent reports whether the invite SMS was actually accepted by the
	// provider. When false, the invite is still recorded; show the link.
	SMSSent bool `json:"sms_sent"`
}

type ListMyReferralsRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" query:"pagination"`
}

// ReferralItem is one row in the "who joined via my link" list.
type ReferralItem struct {
	ID          string     `json:"id"`
	Name        string     `json:"name,omitempty"`
	PhoneMasked string     `json:"phone_masked,omitempty"`
	Source      string     `json:"source"`
	Status      string     `json:"status"`
	Rewarded    bool       `json:"rewarded"`
	InvitedAt   time.Time  `json:"invited_at"`
	JoinedAt    *time.Time `json:"joined_at,omitempty"`
}

type ListMyReferralsResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Items      []ReferralItem                 `json:"items"`
	// Summary mirrors the progress block from GetMyLink for convenience.
	TotalJoined     int `json:"total_joined"`
	PendingProgress int `json:"pending_progress"`
	RewardThreshold int `json:"reward_threshold"`
	NextRewardIn    int `json:"next_reward_in"`
}

// RewardItem is one earned discount code.
type RewardItem struct {
	ID           string    `json:"id"`
	DiscountCode string    `json:"discount_code"`
	CreatedAt    time.Time `json:"created_at"`
}

type ListMyRewardsResponse struct {
	Items []RewardItem `json:"items"`
}

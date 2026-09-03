package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// ReferralStatus is the lifecycle of a single invitee.
type ReferralStatus string

const (
	// ReferralStatusInvited means an SMS invite was sent but the person has not
	// registered yet.
	ReferralStatusInvited ReferralStatus = "invited"
	// ReferralStatusJoined means the invitee completed OTP registration.
	ReferralStatusJoined ReferralStatus = "joined"
)

// ReferralSource records how the invitee reached us.
type ReferralSource string

const (
	ReferralSourceLink ReferralSource = "link"
	ReferralSourceSMS  ReferralSource = "sms"
)

// ReferralAccount is the one-per-user record that owns the permanent invite
// code and tracks progress toward rewards.
type ReferralAccount struct {
	ID     string `bson:"_id" json:"id"`
	UserID string `bson:"user_id" json:"user_id"`
	// Code is the permanent, human-readable invite code (e.g. STV-8K3QP).
	Code string `bson:"code" json:"code"`
	// TotalJoined counts invitees who finished registration.
	TotalJoined int `bson:"total_joined" json:"total_joined"`
	// RewardedCount is how many joined invitees have already been consumed into
	// a reward batch. TotalJoined - RewardedCount is the pending progress.
	RewardedCount int       `bson:"rewarded_count" json:"rewarded_count"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

func (a *ReferralAccount) BeforeCreate() {
	a.ID = ulid.Make().String()
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
}

// PendingProgress is joined invitees not yet turned into a reward.
func (a ReferralAccount) PendingProgress() int {
	p := a.TotalJoined - a.RewardedCount
	if p < 0 {
		return 0
	}
	return p
}

// Referral is one invitee: either an outstanding SMS invite or a completed join.
type Referral struct {
	ID             string         `bson:"_id" json:"id"`
	ReferrerUserID string         `bson:"referrer_user_id" json:"referrer_user_id"`
	Code           string         `bson:"code" json:"code"`
	Source         ReferralSource `bson:"source" json:"source"`
	Status         ReferralStatus `bson:"status" json:"status"`
	// InvitedPhone is set only when the invite was sent by SMS.
	InvitedPhone string `bson:"invited_phone,omitempty" json:"invited_phone,omitempty"`
	// ReferredUserID is filled once the invitee registers.
	ReferredUserID string `bson:"referred_user_id,omitempty" json:"referred_user_id,omitempty"`
	// RewardBatchID links this referral to the ReferralReward that consumed it;
	// empty until it counts toward a granted reward.
	RewardBatchID string     `bson:"reward_batch_id,omitempty" json:"reward_batch_id,omitempty"`
	InvitedAt     time.Time  `bson:"invited_at" json:"invited_at"`
	JoinedAt      *time.Time `bson:"joined_at,omitempty" json:"joined_at,omitempty"`
	CreatedAt     time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `bson:"updated_at" json:"updated_at"`
}

func (r *Referral) BeforeCreate() {
	r.ID = ulid.Make().String()
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	if r.InvitedAt.IsZero() {
		r.InvitedAt = now
	}
	if r.Status == "" {
		r.Status = ReferralStatusInvited
	}
}

// ReferralReward is a discount code granted to a referrer after every N joins.
type ReferralReward struct {
	ID             string    `bson:"_id" json:"id"`
	ReferrerUserID string    `bson:"referrer_user_id" json:"referrer_user_id"`
	DiscountID     string    `bson:"discount_id" json:"discount_id"`
	DiscountCode   string    `bson:"discount_code" json:"discount_code"`
	ReferralIDs    []string  `bson:"referral_ids" json:"referral_ids"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
}

func (r *ReferralReward) BeforeCreate() {
	r.ID = ulid.Make().String()
	r.CreatedAt = time.Now()
}

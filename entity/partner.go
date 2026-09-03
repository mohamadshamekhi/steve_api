package entity

import (
	"time"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/oklog/ulid/v2"
)

// PartnerStatus is the lifecycle state of a partner (colleague-sales) record.
type PartnerStatus string

const (
	PartnerStatusPending   PartnerStatus = "pending"
	PartnerStatusApproved  PartnerStatus = "approved"
	PartnerStatusRejected  PartnerStatus = "rejected"
	PartnerStatusSuspended PartnerStatus = "suspended"
)

// Partner is a storefront user who has been granted access to the colleague
// (B2B) panel. One record per user; the user account itself is unchanged.
type Partner struct {
	ID string `bson:"_id" json:"id"`
	// UserID is the storefront user who submitted the application.
	UserID string `bson:"user_id" json:"user_id"`
	// BackofficeUserID is the restricted backoffice account provisioned for the
	// partner on approval; empty until then. The partner logs into the admin
	// panel with this account and the `partner` role.
	BackofficeUserID string        `bson:"backoffice_user_id,omitempty" json:"backoffice_user_id,omitempty"`
	Status           PartnerStatus `bson:"status" json:"status"`
	DiscountPercent  uint          `bson:"discount_percent" json:"discount_percent"`
	Tier             string        `bson:"tier,omitempty" json:"tier,omitempty"`

	CompanyName   string `bson:"company_name" json:"company_name"`
	NationalID    string `bson:"national_id" json:"national_id"`
	BusinessField string `bson:"business_field" json:"business_field"`
	City          string `bson:"city" json:"city"`

	RejectionReason string     `bson:"rejection_reason,omitempty" json:"rejection_reason,omitempty"`
	ReviewedBy      string     `bson:"reviewed_by,omitempty" json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time `bson:"reviewed_at,omitempty" json:"reviewed_at,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (p *Partner) BeforeCreate() {
	p.ID = ulid.Make().String()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.Status == "" {
		p.Status = PartnerStatusPending
	}
}

func (p Partner) IsApproved() bool {
	return p.Status == PartnerStatusApproved
}

// CanTransitionPartnerStatus mirrors the state machine drawn in the feature
// design: pending → approved/rejected, rejected → pending (re-apply),
// approved ↔ suspended.
func CanTransitionPartnerStatus(from, to PartnerStatus) bool {
	allowed := map[PartnerStatus][]PartnerStatus{
		PartnerStatusPending:   {PartnerStatusApproved, PartnerStatusRejected},
		PartnerStatusRejected:  {PartnerStatusPending},
		PartnerStatusApproved:  {PartnerStatusSuspended},
		PartnerStatusSuspended: {PartnerStatusApproved},
	}
	for _, s := range allowed[from] {
		if s == to {
			return true
		}
	}
	return false
}

func (p Partner) GuardTransition(to PartnerStatus) error {
	const op = "entity.Partner.GuardTransition"
	if !CanTransitionPartnerStatus(p.Status, to) {
		return richerror.New(op).
			WithKind(richerror.KindBadRequest).
			WithMessageKey("partner.invalid_status_transition")
	}
	return nil
}

package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// Supplier is a shop or person the store buys goods from — the "مغازه بغل"
// an out-of-stock item gets sourced from, or a regular wholesaler.
type Supplier struct {
	ID          string    `bson:"_id" json:"id"`
	Name        string    `bson:"name" json:"name"`
	ShopName    string    `bson:"shop_name,omitempty" json:"shop_name,omitempty"`
	PhoneNumber string    `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	Address     string    `bson:"address,omitempty" json:"address,omitempty"`
	Note        string    `bson:"note,omitempty" json:"note,omitempty"`
	IsActive    bool      `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt   time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func (s *Supplier) BeforeCreate() {
	s.ID = ulid.Make().String()
	now := time.Now()
	s.IsActive = true
	s.CreatedAt = now
	s.UpdatedAt = now
}

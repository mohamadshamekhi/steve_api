package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Shipping struct {
	ID          string    `bson:"_id" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Cost        uint64    `bson:"cost" json:"cost"`
	IsActive    bool      `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

func (sh *Shipping) BeforeCreate() {
	sh.ID = ulid.Make().String()
	sh.CreatedAt = time.Now()
	sh.UpdatedAt = sh.CreatedAt
}

type Banner struct {
	ID    string `json:"id" bson:"_id"`
	Title string `json:"title" bson:"title"`

	MediaID string `json:"media_id" bson:"media_id"`

	Link string `json:"link" bson:"link"`

	Location string `json:"location" bson:"location"`

	SortOrder int `json:"sort_order" bson:"sort_order"`

	IsActive bool `json:"is_active" bson:"is_active"`

	AltText string `json:"alt_text" bson:"alt_text"`

	StartsAt *time.Time `json:"starts_at,omitempty" bson:"starts_at,omitempty"`

	ExpiresAt *time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (b *Banner) BeforeCreate() {
	b.ID = ulid.Make().String()

	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
}

package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Brand struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description,omitempty"`
	IsFavorite  bool      `bson:"is_favorite,omitempty"`
	MediaID     string    `bson:"media_id,omitempty"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
	DeletedAt   time.Time `bson:"deleted_at,omitempty"`
}

func (b *Brand) BeforeCreate() {
	b.ID = ulid.Make().String()
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
}

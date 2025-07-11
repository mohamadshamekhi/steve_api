package entity

import (
	"github.com/oklog/ulid/v2"
	"time"
)

type Article struct {
	ID          string    `bson:"_id"`
	Title       string    `bson:"title"`
	Slug        string    `bson:"slug"`
	Body        string    `bson:"body"`
	IsPublished bool      `bson:"is_published"`
	MediaID     string    `bson:"media_id"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
	DeletedAt   time.Time `bson:"deleted_at"`
}

func (a *Article) BeforeCreate() {
	a.ID = ulid.Make().String()
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
}

package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Article struct {
	ID    string `bson:"_id"`
	Title string `bson:"title"`
	Slug  string `bson:"slug"`
	// PreviousSlugs keeps every slug the article had before a rename, so the
	// storefront can 301-redirect old links.
	PreviousSlugs   []string  `bson:"previous_slugs,omitempty"`
	Body            string    `bson:"body"`
	IsPublished     bool      `bson:"is_published"`
	Tags            []string  `bson:"tags"`
	MetaDescription string    `bson:"meta_description"`
	BackOfficeID    string    `bson:"back_office_id"`
	MediaID         string    `bson:"media_id"`
	Alt             string    `bson:"alt"`
	ViewCount       int64     `bson:"view_count"`
	CreatedAt       time.Time `bson:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at"`
	DeletedAt       time.Time `bson:"deleted_at,omitempty"`
}

func (a *Article) BeforeCreate() {
	a.ID = ulid.Make().String()
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
}

func (a Article) GetTitle() string {
	return a.Title
}
func (a Article) GetID() string {
	return a.ID
}

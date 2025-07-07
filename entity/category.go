package entity

import (
	"github.com/oklog/ulid/v2"
	"time"
)

type Category struct {
	ID          string     `bson:"_id"`
	Name        string     `bson:"name"`
	Slug        string     `bson:"slug"`
	Description string     `bson:"description,omitempty"`
	MediaID     string     `bson:"media_id"`
	ParentID    string     `bson:"parent_id,omitempty"`
	Ancestors   []Ancestor `bson:"ancestors,omitempty"`
	Level       int        `bson:"level"`
	CreatedAt   time.Time  `bson:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at"`
	DeletedAt   time.Time  `bson:"deleted_at,omitempty"`
}

type Ancestor struct {
	ID    string `bson:"_id"`
	Name  string `bson:"name"`
	Slug  string `bson:"slug"`
	Level int    `bson:"level"`
}

func (c *Category) BeforeCreate() {
	c.ID = ulid.Make().String()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
}

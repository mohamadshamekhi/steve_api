package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Comment struct {
	ID string `bson:"_id"`

	Body string `bson:"body"`

	CommentableID   string `bson:"commentable_id"`
	CommentableType string `bson:"commentable_type"` //  (مثلاً "product" || "article")

	CommenterID   string `bson:"commenter_id"`
	CommenterType string `bson:"commenter_type"` //  ("user" || "backoffice")
	CommenterName string `bson:"commenter_name"`

	ParentID string `bson:"parent_id,omitempty"`

	Ancestors []AncestorComment `bson:"ancestors,omitempty"`

	Level int `bson:"level"`

	IsApproved bool      `bson:"is_approved"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
	DeletedAt  time.Time `bson:"deleted_at,omitempty"`
}
type AncestorComment struct {
	ID            string `bson:"_id" json:"id"`
	CommenterID   string `bson:"commenter_id" json:"commenter_id"`
	CommenterType string `bson:"commenter_type" json:"commenter_type"`
	CommenterName string `bson:"commenter_name" json:"commenter_name"`
	Body          string `bson:"body" json:"body"`
	Level         int    `bson:"level" json:"level"`
}

func (c *Comment) BeforeCreate() {
	c.ID = ulid.Make().String()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.IsApproved = false
}

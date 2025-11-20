package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type ContactUs struct {
	ID          string    `bson:"_id"`
	Email       string    `bson:"email"`
	PhoneNumber string    `bson:"phone_number"`
	Subject     string    `bson:"subject"`
	Description string    `bson:"description"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
	DeletedAt   time.Time `bson:"deleted_at,omitempty"`
}

func (c *ContactUs) BeforeCreate() {
	c.ID = ulid.Make().String()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
}

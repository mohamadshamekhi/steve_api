package entity

import (
	"github.com/oklog/ulid/v2"
	"time"
)

type User struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	PhoneNumber string    `bson:"phone_number"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func (u *User) BeforeCreate() {
	u.ID = ulid.Make().String()
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
}

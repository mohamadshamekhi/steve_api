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
	Addresses   []Address `bson:"addresses,omitempty"`
}

type Address struct {
	ID       string `bson:"id"`
	Title    string `bson:"title"`
	Province string `bson:"province"`
	City     string `bson:"city"`
	Street   string `bson:"street"`
	Postcode string `bson:"postcode"`
}

func (u *User) BeforeCreate() {
	u.ID = ulid.Make().String()
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
}

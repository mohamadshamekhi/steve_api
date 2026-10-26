package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type User struct {
	ID          string     `bson:"_id"`
	Name        string     `bson:"name"`
	PhoneNumber string     `bson:"phone_number"`
	BirthDate   *time.Time `bson:"birth_date,omitempty"`
	CreatedAt   time.Time  `bson:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at"`
	Addresses   []Address  `bson:"addresses,omitempty"`
}

type Address struct {
	ID       string `bson:"id" json:"id"`
	Title    string `bson:"title" json:"title"`
	Province string `bson:"province" json:"province"`
	City     string `bson:"city" json:"city"`
	Street   string `bson:"street" json:"street"`
	Postcode string `bson:"postcode" json:"postcode"`
}

func (u *User) BeforeCreate() {
	u.ID = ulid.Make().String()
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
}
func (u User) GetID() string {
	return u.ID
}
func (u User) GetName() string {
	return u.Name
}

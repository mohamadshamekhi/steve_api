package entity

import (
	"github.com/oklog/ulid/v2"
	"time"
)

type BackOffice struct {
	ID           string    `bson:"_id" json:"id"`
	Name         string    `bson:"name" json:"name"`
	PhoneNumber  string    `bson:"phone_number" json:"phone_number"`
	NationalCode string    `bson:"national_code" json:"national_code"`
	Address      string    `bson:"address" json:"address"`
	ZipCode      string    `bson:"zip_code" json:"zip_code"`
	AvatarID     int64     `bson:"avatar_id" json:"avatar_id"`
	IsActive     bool      `bson:"is_active" json:"is_active"`
	isSuperAdmin bool      `bson:"is_super_admin"`
	Roles        []string  `bson:"roles,omitempty" json:"roles,omitempty"`
	Permissions  []string  `bson:"permissions,omitempty" json:"permissions,omitempty"`
	CreatedAt    time.Time `bson:"created_at" json:"-"`
	UpdatedAt    time.Time `bson:"updated_at" json:"-"`
}

func (b BackOffice) IsSuperAdmin() bool {
	return b.isSuperAdmin
}

func (b *BackOffice) BeforeCreate() {
	b.ID = ulid.Make().String()
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
}

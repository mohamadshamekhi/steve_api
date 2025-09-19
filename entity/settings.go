package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Shipping struct {
	ID          string    `bson:"_id" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Cost        uint64    `bson:"cost" json:"cost"`
	IsActive    bool      `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

func (sh *Shipping) BeforeCreate() {
	sh.ID = ulid.Make().String()
	sh.CreatedAt = time.Now()
	sh.UpdatedAt = sh.CreatedAt
}

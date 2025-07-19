package entity

import (
	"github.com/oklog/ulid/v2"
	"time"
)

type AttributeDefinition struct {
	ID   string `bson:"_id"`
	Key  string `bson:"key"`            // example : "cpu_processor", "main_camera_resolution"
	Name string `bson:"name"`           //example : "پردازنده", "رزولوشن دوربین اصلی"
	Unit string `bson:"unit,omitempty"` //example "مگاپیکسل", "گیگابایت", "اینچ"
}

type AttributeSet struct {
	ID         string                `bson:"_id"`
	Name       string                `bson:"name"` //example:"مشخصات فنی موبایل", "مشخصات کتاب"
	Attributes []AttributeDefinition `bson:"attributes"`
	CreatedAt  time.Time             `bson:"created_at"`
	UpdatedAt  time.Time             `bson:"updated_at"`
}

func (a *AttributeSet) BeforeCreate() {
	a.ID = ulid.Make().String()
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
}

func (a *AttributeDefinition) BeforeCreate() {
	a.ID = ulid.Make().String()
}

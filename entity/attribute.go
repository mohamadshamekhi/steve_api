package entity

import (
	"github.com/oklog/ulid/v2"
	"time"
)

type AttributeDefinition struct {
	ID   string `bson:"_id" json:"id"`
	Key  string `bson:"key" json:"key"`             // example : "cpu_processor", "main_camera_resolution"
	Name string `bson:"name" json:"name"`           //example : "پردازنده", "رزولوشن دوربین اصلی"
	Unit string `bson:"unit,omitempty" json:"unit"` //example "مگاپیکسل", "گیگابایت", "اینچ"
}

type AttributeSet struct {
	ID         string                `bson:"_id" json:"id"`
	Name       string                `bson:"name" json:"name"` //example:"مشخصات فنی موبایل", "مشخصات کتاب"
	Attributes []AttributeDefinition `bson:"attributes" json:"attributes"`
	CreatedAt  time.Time             `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time             `bson:"updated_at" json:"updated_at"`
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

package entity

import (
	"github.com/oklog/ulid/v2"
)

type VariantAttribute struct {
	Key   string `bson:"key" json:"key"`               // ex: "color"
	Value string `bson:"value" json:"value"`           // ex: "red"
	Label string `bson:"label,omitempty" json:"label"` // ex:"قزمز",
	Meta  string `bson:"meta,omitempty" json:"meta"`
}

type Variant struct {
	ID               string             `bson:"id" json:"id"`
	Price            uint64             `bson:"price" json:"price"`
	Quantity         uint64             `bson:"quantity" json:"quantity"`
	VariantAttribute []VariantAttribute `bson:"variant_attribute,omitempty" json:"variant_attribute,omitempty"`
	Warranty         string             `bson:"warranty" json:"warranty"`
}

func (v *Variant) BeforeCreate() {
	v.ID = ulid.Make().String()
}

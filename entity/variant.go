package entity

import (
	"github.com/oklog/ulid/v2"
)

type VariantAttribute struct {
	Key      string `bson:"key"`             // ex: "color"
	Value    string `bson:"value"`           // ex: "red"
	Label    string `bson:"label,omitempty"` // ex:"قزمز",
	HexColor string `bson:"hex_color,omitempty"`
}

type Variant struct {
	ID               string             `bson:"id"`
	Price            uint64             `bson:"price"`
	Quantity         uint64             `bson:"quantity"`
	VariantAttribute []VariantAttribute `bson:"variant_attribute,omitempty"`
}

func (v *Variant) BeforeCreate() {
	v.ID = ulid.Make().String()
}

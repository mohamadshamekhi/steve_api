package entity

type VariantAttribute struct {
	Key      string `bson:"key"`             // ex: "color"
	Value    string `bson:"value"`           // ex: "red"
	Label    string `bson:"label,omitempty"` // ex:"قزمز",
	HexColor string `bson:"hex_color,omitempty"`
}

type Variant struct {
	ID               string             `bson:"id"`
	Price            int64              `bson:"price"`
	Quantity         int                `bson:"quantity"`
	VariantAttribute []VariantAttribute `bson:"variant_attribute,omitempty"`
}

package productparam

type AddVariantRequest struct {
	ProductID string `bson:"product_id"`
	Price     uint64 `bson:"price"`
	Warranty  string `bson:"warranty"`
	Quantity  uint64 `bson:"quantity"`
	Key       string `bson:"key"`
	Value     string `bson:"value"`
	Label     string `bson:"label,omitempty"`
	Meta      string `bson:"meta,omitempty"`
}

type AddVariantResponse struct {
}

type DeleteVariantRequest struct {
	ProductID string `bson:"product_id"`
	VariantID string `bson:"variant_id"`
}

type DeleteVariantResponse struct {
}

type UpdateVariantRequest struct {
	ProductID string `bson:"product_id"`
	VariantID string `bson:"variant_id"`

	Price    *uint64 `bson:"price"`
	Warranty *string `bson:"warranty"`
	Quantity *uint64 `bson:"quantity"`
	Key      *string `bson:"key"`
	Value    *string `bson:"value"`
	Label    *string `bson:"label,omitempty"`
	Meta     *string `bson:"meta,omitempty"`
}

type UpdateVariantResponse struct {
}

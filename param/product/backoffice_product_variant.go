package productparam

type AddVariantRequest struct {
	ProductID string `json:"product_id" bson:"product_id"`
	Price     uint64 `json:"price" bson:"price"`
	Warranty  string `json:"warranty" bson:"warranty"`
	Quantity  uint64 `json:"quantity" bson:"quantity"`
	Key       string `json:"key" bson:"key"`
	Value     string `json:"value" bson:"value"`
	Label     string `json:"label,omitempty" bson:"label,omitempty"`
	Meta      string `json:"meta,omitempty" bson:"meta,omitempty"`
}

type AddVariantResponse struct {
}

type DeleteVariantRequest struct {
	ProductID string `json:"product_id" bson:"product_id"`
	VariantID string `json:"variant_id" bson:"variant_id"`
}

type DeleteVariantResponse struct {
}

type UpdateVariantRequest struct {
	ProductID string `json:"product_id" bson:"product_id"`
	VariantID string `json:"variant_id" bson:"variant_id"`

	Price    *uint64 `json:"price" bson:"price"`
	Warranty *string `json:"warranty" bson:"warranty"`
	Quantity *uint64 `json:"quantity" bson:"quantity"`
	Key      *string `json:"key" bson:"key"`
	Value    *string `json:"value" bson:"value"`
	Label    *string `json:"label,omitempty" bson:"label,omitempty"`
	Meta     *string `json:"meta,omitempty" bson:"meta,omitempty"`
}

type UpdateVariantResponse struct {
}

// BulkUpdateVariantsRequest reuses UpdateVariantRequest per item: same fields
// as the single update-variant body, plus product_id/variant_id in the item
// itself since there's no per-product URL for a bulk request.
type BulkUpdateVariantsRequest struct {
	Items []UpdateVariantRequest `json:"items"`
}

type BulkUpdateVariantItemResult struct {
	VariantID string `json:"variant_id"`
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
}

type BulkUpdateVariantsResponse struct {
	Results []BulkUpdateVariantItemResult `json:"results"`
}

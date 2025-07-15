package variantparam

type UpdateVariantRequest struct {
	ID       string  `json:"id"`
	Price    *uint64 `json:"price"`
	Quantity *uint64 `json:"quantity"`
	Key      *string `json:"key"`
	Value    *string `json:"value"`
	Label    *string `json:"label,omitempty"`
	HexColor *string `json:"hex_color,omitempty"`
}

type UpdateVariantResponse struct {
}

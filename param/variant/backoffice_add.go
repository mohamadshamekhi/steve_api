package variantparam

type AddVariantRequest struct {
	Price    uint64 `json:"price"`
	Quantity uint64 `json:"quantity"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	Label    string `json:"label,omitempty"`
	HexColor string `json:"hex_color,omitempty"`
}

type AddVariantResponse struct {
}

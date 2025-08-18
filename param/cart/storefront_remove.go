package cartparam

type RemoveFromCartRequest struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
}

type RemoveFromCartResponse struct {
}

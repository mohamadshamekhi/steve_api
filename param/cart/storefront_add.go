package cartparam

type AddItemToCartRequest struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
}

type AddItemToCartResponse struct {
}

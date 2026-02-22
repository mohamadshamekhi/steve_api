package cartparam

type RemoveFromCartBackofficeRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
}

type RemoveFromCartBackofficeResponse struct {
}

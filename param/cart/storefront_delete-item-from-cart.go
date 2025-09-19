package cartparam

type DeleteItemFromCartRequest struct {
	UserID    string `json:"user_id"`
	VariantID string `json:"variant_id"`
}

type DeleteItemFromCartResponse struct {
}

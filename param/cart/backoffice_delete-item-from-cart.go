package cartparam

type DeleteItemFromCartBackofficeRequest struct {
	UserID    string `json:"user_id"`
	VariantID string `json:"variant_id"`
}

type DeleteItemFromCartBackofficeResponse struct {
}

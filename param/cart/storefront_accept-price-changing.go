package cartparam

type AcceptPriceChangeRequest struct {
	VariantID string `json:"variant_id"`
	UserID    string `json:"user_id"`
}
type AcceptPriceChangeResponse struct{}

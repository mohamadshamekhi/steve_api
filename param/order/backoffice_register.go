package orderparam

type RegisterOrderBackofficeRequest struct {
	UserID       string `json:"user_id"`
	DiscountCode string `json:"discount_code"`
}

type RegisterOrderBackofficeResponse struct {
	OrderID string `json:"order_id"`
}

package orderparam

type CompletedOrderBackofficeRequest struct {
	OrderID     string `json:"order_id"`
	Description string `json:"description"`
}

type CompletedOrderBackofficeResponse struct {
}

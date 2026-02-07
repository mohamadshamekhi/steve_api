package paymentparam

type GoToGatewayRequest struct {
	Amount      uint   `json:"amount"`
	Description string `json:"description"`
	OrderID     string `json:"order_id"`
	UserID      string `json:"user_id"`
}
type GoToGatewayResponse struct {
	GatewayTransactionID string `json:"gateway_transaction_id"`
	PaymentURL           string `json:"payment_url"`
}

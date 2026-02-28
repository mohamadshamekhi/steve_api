package paymentparam

type CreatePaymentForBackofficePanelRequest struct {
	OrderID     string `json:"order_id"`
	Description string `json:"description"`
}

type CreatePaymentForBackofficePanelResponse struct {
}

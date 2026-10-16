package orderparam

import "github.com/abolfazlnorzad/steve/entity"

type CompletedOrderBackofficeRequest struct {
	OrderID     string `json:"order_id"`
	Description string `json:"description"`
	// PaymentMethod is how the customer paid at the counter; it is recorded
	// in the cashbox. Empty defaults to cash.
	PaymentMethod entity.PaymentMethod `json:"payment_method"`
}

type CompletedOrderBackofficeResponse struct {
}

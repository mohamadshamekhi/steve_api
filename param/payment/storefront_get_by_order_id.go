package paymentparam

import "github.com/abolfazlnorzad/steve/entity"

type GetPaymentsByOrderIDRequest struct {
	OrderID string `json:"order_id"`
}
type GetPaymentsByOrderIDResponse struct {
	Payments []entity.Payment `json:"payments"`
}

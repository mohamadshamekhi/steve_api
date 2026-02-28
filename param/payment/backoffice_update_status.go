package paymentparam

import "github.com/abolfazlnorzad/steve/entity"

type UpdateStatusRequest struct {
	PaymentID string               `json:"payment_id"`
	Status    entity.PaymentStatus `json:"status"`
}

type UpdateStatusResponse struct {
}

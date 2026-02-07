package paymentparam

import "github.com/abolfazlnorzad/steve/entity"

type VerifyPaymentRequest struct {
	GatewayTransactionID string `json:"gateway_transaction_id"`
}

type VerifyPaymentResponse struct {
	Status       entity.PaymentStatus `json:"status"`
	TrackingCode string               `json:"tracking_code"`
	Message      string               `json:"message"`
}

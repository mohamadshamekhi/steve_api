package orderparam

import "github.com/abolfazlnorzad/steve/entity"

type UpdateStatusOrderRequest struct {
	OrderID string             `json:"order_id"`
	Status  entity.OrderStatus `json:"status"`
	// TrackingCode is required when Status is shipped.
	TrackingCode string `json:"tracking_code"`
}

type UpdateStatusOrderResponse struct {
}

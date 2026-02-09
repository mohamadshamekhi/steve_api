package orderparam

import "github.com/abolfazlnorzad/steve/entity"

type UpdateStatusOrderRequest struct {
	OrderID string             `json:"order_id"`
	Status  entity.OrderStatus `json:"status"`
}

type UpdateStatusOrderResponse struct {
}

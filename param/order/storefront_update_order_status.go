package orderparam

import "github.com/abolfazlnorzad/steve/entity"

type UpdateStatusRequest struct {
	OrderID string             `json:"order_id"`
	Status  entity.OrderStatus `json:"status"`
}
type UpdateStatusResponse struct {
}

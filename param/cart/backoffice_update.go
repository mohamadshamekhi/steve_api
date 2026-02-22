package cartparam

import "github.com/abolfazlnorzad/steve/entity"

type UpdateCartStatusRequest struct {
	CartID string            `json:"cart_Id"`
	Status entity.CartStatus `json:"status"`
}

type UpdateCartStatusResponse struct {
}

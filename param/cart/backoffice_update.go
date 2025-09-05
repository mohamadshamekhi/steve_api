package cartparam

import "github.com/abolfazlnorzad/steve/entity"

type UpdateCartStatusRequest struct {
	CartID string
	Status entity.CartStatus
}

type UpdateCartStatusResponse struct {
}

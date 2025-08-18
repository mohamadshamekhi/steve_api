package cartparam

import "github.com/abolfazlnorzad/steve/entity"

type GetMyActiveCartRequest struct {
}

type GetMyActiveCartResponse struct {
	Cart entity.Cart `json:"cart"`
}

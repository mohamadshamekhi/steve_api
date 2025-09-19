package sellparam

import "github.com/abolfazlnorzad/steve/entity"

type UpdateSellRequest struct {
	ID     string            `json:"id"`
	Status entity.SellStatus `json:"status"`
}

type UpdateSellResponse struct {
}

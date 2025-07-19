package productparam

import "github.com/abolfazlnorzad/steve/entity"

type GetProductRequest struct {
	ID string `json:"id"`
}
type GetProductResponse struct {
	Product entity.Product `json:"product"`
}

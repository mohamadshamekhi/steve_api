package productparam

import "github.com/abolfazlnorzad/steve/entity"

type AddProductAttributeRequest struct {
	Attrs     []entity.Specification `json:"attrs"`
	ProductID string                 `json:"product_id"`
}

type AddProductAttributeResponse struct {
}

type UpdateSingleAttributeRequest struct {
	ProductID string `json:"product_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

type UpdateSingleAttributeResponse struct {
}

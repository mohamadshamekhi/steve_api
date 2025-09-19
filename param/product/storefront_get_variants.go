package productparam

import "github.com/abolfazlnorzad/steve/entity"

type GetVariantsByIDsRequest struct {
	IDs []string `json:"ids"`
}

type GetVariantsByIDsResponse struct {
	Variants map[string]entity.Variant `json:"variants"`
}

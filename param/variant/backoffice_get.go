package variantparam

import "github.com/abolfazlnorzad/steve/entity"

type GetVariantRequest struct {
	ID string `json:"id"`
}

type GetVariantResponse struct {
	Variant entity.Variant `json:"variant"`
}

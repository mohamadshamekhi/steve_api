package attributeparam

import "github.com/abolfazlnorzad/steve/entity"

type GetAttributeRequest struct {
	ID string `json:"id"`
}

type GetAttributeResponse struct {
	Attr entity.AttributeSet `json:"attr"`
}

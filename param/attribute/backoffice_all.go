package attributeparam

import (
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type AllAttributeRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type AllAttributeResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Attributes []entity.AttributeSet          `json:"attributes"`
}

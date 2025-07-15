package attributedefinitionparam

import (
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type AllAttributeDefinitionRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type AllAttributeDefinitionResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Attributes []entity.AttributeDefinition   `json:"attributes"`
}

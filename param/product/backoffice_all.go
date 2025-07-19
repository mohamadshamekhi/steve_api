package productparam

import (
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllProductsRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllProductsResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Products   []entity.Product               `json:"products"`
}

package categoryparam

import (
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllCategoriesParamRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllCategoriesParamResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Brands     []CategoryInfo                 `json:"categories"`
}

type CategoryInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	URL         string            `json:"url"`
	Ancestors   []entity.Ancestor `json:"ancestors"`
}

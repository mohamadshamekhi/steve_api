package categoryparam

import (
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllCategoriesRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllCategoriesResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Categories []CategoryInfo                 `json:"categories"`
}

type CategoryInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	URL         string            `json:"url"`
	Level       int               `json:"level"`
	Slug        string            `json:"slug"`
	Ancestors   []entity.Ancestor `json:"ancestors"`
}

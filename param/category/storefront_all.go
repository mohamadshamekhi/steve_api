package categoryparam

import commonparam "github.com/abolfazlnorzad/steve/param/common"

type GetAllStorefrontCategoriesRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllStorefrontCategoriesResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Categories []CategoryInfo                 `json:"categories"`
}

package brandparam

import (
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllBrandsParamRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllBrandsParamResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Brands     []BrandInfo                    `json:"brands"`
}

type BrandInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	IsFavorite  bool   `json:"is_favorite"`
}

package brandparam

import commonparam "github.com/abolfazlnorzad/steve/param/common"

type GetAllStorefrontBrandsParamRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllStorefrontBrandsParamResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Brands     []BrandInfo                    `json:"brands"`
}

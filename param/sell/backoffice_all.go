package sellparam

import commonparam "github.com/abolfazlnorzad/steve/param/common"

type GetAllSellParamRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllSellParamResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Sells      []SellInfo                     `json:"sells"`
}

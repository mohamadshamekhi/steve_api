package repairparam

import commonparam "github.com/abolfazlnorzad/steve/param/common"

type GetAllRepairParamRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllRepairParamResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Repairs    []RepairInfo                   `json:"repairs"`
}

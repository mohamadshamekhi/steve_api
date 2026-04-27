package backofficeparam

import (
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type AllBackofficeRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}

type AllBackofficeResponse struct {
	Backoffice []entity.BackOffice            `json:"backoffice"`
	Pagination commonparam.PaginationResponse `json:"pagination"`
}

type GetSuperUserBackofficeRequest struct{}

type GetSuperUserBackofficeResponse struct {
	Backoffice []entity.BackOffice `json:"backoffice"`
}

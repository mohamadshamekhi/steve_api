package contactusparam

import (
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllContactRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllContactResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Contacts   []entity.ContactUs             `json:"brands"`
}

package settingparam

import (
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type AddShippingRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Cost        uint64 `json:"cost"`
	IsActive    bool   `json:"is_active"`
}

type AddShippingResponse struct {
}

type UpdateShippingRequest struct {
	ShippingID  string  `json:"shipping_id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Cost        *uint64 `json:"cost"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateShippingResponse struct {
}

type FindByIDShippingRequest struct {
	ShippingID string `json:"shipping_id"`
}

type FindByIDShippingResponse struct {
	Shipping entity.Shipping `json:"shipping"`
}

type GetAllShippingRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     FilterShipping                `json:"filter" ,query:"filter"`
}

type FilterShipping struct {
	IsActive *bool `json:"is_active" query:"is_active"`
}

type GetAllShippingResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Shipping   []entity.Shipping              `json:"shipping" bson:"shipping"`
}

type GetActiveShippingRequest struct {
}

type GetActiveShippingResponse struct {
	Shipping []entity.Shipping `json:"shipping"`
}

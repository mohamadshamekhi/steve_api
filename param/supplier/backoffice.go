package supplierparam

import (
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type CreateSupplierRequest struct {
	Name        string `json:"name"`
	ShopName    string `json:"shop_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Note        string `json:"note"`
}

type CreateSupplierResponse struct {
	Supplier entity.Supplier `json:"supplier"`
}

type UpdateSupplierRequest struct {
	ID          string  `json:"-"`
	Name        *string `json:"name"`
	ShopName    *string `json:"shop_name"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
	Note        *string `json:"note"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateSupplierResponse struct {
	Supplier entity.Supplier `json:"supplier"`
}

type ListSuppliersRequest struct {
	Search     string
	IsActive   *bool
	Pagination commonparam.PaginationRequest
}

type ListSuppliersResponse struct {
	Suppliers  []entity.Supplier              `json:"suppliers"`
	Pagination commonparam.PaginationResponse `json:"pagination"`
}

type GetSupplierResponse struct {
	Supplier entity.Supplier `json:"supplier"`
}

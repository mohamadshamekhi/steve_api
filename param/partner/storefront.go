package partnerparam

import (
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
)

// ApplyRequest is the colleague-partnership application form. Economic code and
// business-licence upload are intentionally not collected.
type ApplyRequest struct {
	CompanyName   string `json:"company_name"`
	NationalID    string `json:"national_id"`
	BusinessField string `json:"business_field"`
	City          string `json:"city"`
}

type ApplyResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// MeResponse powers the partner panel's bootstrap call.
type MeResponse struct {
	HasPartnership  bool   `json:"has_partnership"`
	Status          string `json:"status,omitempty"`
	DiscountPercent uint   `json:"discount_percent,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	CompanyName     string `json:"company_name,omitempty"`
}

type CatalogListRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}

type CatalogVariant struct {
	ID             string `json:"id"`
	Price          uint64 `json:"price"`
	PartnerPrice   uint64 `json:"partner_price"`
	DiscountAmount uint64 `json:"discount_amount"`
	Quantity       uint64 `json:"quantity"`
	Warranty       string `json:"warranty"`
}

type CatalogProduct struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Slug            string           `json:"slug"`
	URL             string           `json:"url"`
	PartnerEligible bool             `json:"partner_eligible"`
	AppliedPercent  uint             `json:"applied_percent"`
	Variants        []CatalogVariant `json:"variants"`
}

type CatalogListResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Products   []CatalogProduct               `json:"products"`
}

type MyOrdersResponse struct {
	Pagination commonparam.PaginationResponse   `json:"pagination"`
	Orders     []orderparam.StoreFrontOrderInfo `json:"orders"`
}

package productparam

import (
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllProductsRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllProductsResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Products   []ProductList                  `json:"products"`
}

type ProductList struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	EnTitle                 string `json:"en_title"`
	Brand                   string `json:"brand"`
	Category                string `json:"category"`
	IsPublished             bool   `json:"is_published"`
	IsSecondHand            bool   `json:"is_second_hand"`
	RequiresFreightShipping bool   `json:"requires_freight_shipping"`
	Slug                    string `json:"slug"`
	ThumbnailURL            string `json:"thumbnail_url,omitempty"`

	Kind          entity.ProductKind    `json:"kind"`
	SalesChannels []entity.SalesChannel `json:"sales_channels"`
	// Price is the fee of a service, so the services list can show it without
	// a second request. Omitted for goods, which have one price per variant.
	Price uint64 `json:"price,omitempty"`
}

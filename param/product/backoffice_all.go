package productparam

import (
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
	ID           string `json:"id"`
	Name         string `json:"name"`
	EnTitle      string `json:"en_title"`
	Brand        string `json:"brand"`
	Category     string `json:"category"`
	IsPublished  bool   `json:"is_published"`
	IsSecondHand bool   `json:"is_second_hand"`
	Slug         string `json:"slug"`
}

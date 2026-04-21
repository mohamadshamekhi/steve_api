package productparam

import commonparam "github.com/abolfazlnorzad/steve/param/common"

type StorefrontGetAllProductRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     StorefrontFilterProduct
}
type StorefrontFilterProduct struct {
	BrandID      *string        `json:"brand_id"`
	CategoryID   *string        `json:"category_id"`
	CategorySlug *string        `json:"category_slug"`
	MinPrice     *uint64        `json:"min_price"`
	MaxPrice     *uint64        `json:"max_price"`
	Color        *string        `json:"color"`
	IsExist      *bool          `json:"is_exist"`
	IsSecondHand *bool          `json:"is_second_hand"`
	Filter       map[string]any `json:"filter"`
	OrderBy      *string        `json:"order_by"`
}

type StorefrontGetAllProductResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Products   []StorefrontProductList        `json:"products"`
}

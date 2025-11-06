package productparam

import (
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetProductsWithAutoDiscountsRequest struct{}

type GetProductsWithAutoDiscountsResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Products   []StorefrontProductList
}

type StorefrontProductList struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	ENTitle      string                  `json:"en_title"`
	Slug         string                  `json:"slug"`
	Gallery      []entity.ProductGallery `json:"gallery"`
	Rating       float64                 `json:"rating"`
	IsPublished  bool                    `json:"is_published"`
	IsSecondHand bool                    `json:"is_second_hand"`
	Discount     entity.Discount         `json:"discount"` // پوینتر برای پشتیبانی از null
	Variant      []entity.Variant        `json:"variant" bson:"variant"`
}

type GetProductWithSpecialCategoryRequest struct {
	CategoryID string
}

type GetProductWithSpecialCategoryResponse struct {
	ProductID string
}

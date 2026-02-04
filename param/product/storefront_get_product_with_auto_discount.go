package productparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetProductsWithAutoDiscountsRequest struct{}

type GetProductsWithAutoDiscountsResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Products   []StorefrontProductList
}

type StorefrontProductList struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	ENTitle        string                  `json:"en_title"`
	Slug           string                  `json:"slug"`
	Gallery        []entity.ProductGallery `json:"gallery"`
	Rating         float64                 `json:"rating"`
	IsPublished    bool                    `json:"is_published"`
	IsSecondHand   bool                    `json:"is_second_hand"`
	Discount       []ProductDiscount       `json:"discount"`
	Variant        []entity.Variant        `json:"variant" bson:"variant"`
	Specifications []Specification         `json:"specifications,omitempty"`
}

type GetProductWithSpecialCategoryRequest struct {
	CategoryID string
}

type GetProductWithSpecialCategoryResponse struct {
	ProductID string
}

type ProductDiscount struct {
	DiscountID    string              `json:"discount_id"`
	DiscountType  string              `json:"discount_type"`
	DiscountValue uint                `json:"discount_value"`
	StartAt       time.Time           `json:"start_at"`
	ExpiresAt     time.Time           `json:"expires_at"`
	IsActive      bool                `json:"is_active"`
	IsSpecial     bool                `json:"is_special"`
	ApplicableTo  entity.ApplicableTo `json:"applicable_to"`
	VariantID     string              `json:"variant_id"`
	ProductID     string              `json:"product_id"`
}

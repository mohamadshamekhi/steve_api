package productparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
)

type StorefrontGetProductRequest struct {
	ID string `json:"id"`
}
type StorefrontGetProductResponse struct {
	Product StorefrontProductInfo `json:"product"`
}

type StorefrontGetProductBySlugRequest struct {
	Slug string `json:"slug"`
}
type StorefrontGetProductBySlugResponse struct {
	Product StorefrontProductInfo `json:"product"`
}

// StorefrontProductInfo mirrors ProductInfo but is safe for unauthenticated
// storefront responses: Variants uses StorefrontVariant, which omits
// internal cost/supplier fields.
type StorefrontProductInfo struct {
	ID                      string                     `json:"_id"`
	Name                    string                     `json:"name"`
	ENTitle                 string                     `json:"en_title"`
	Slug                    string                     `json:"slug"`
	Brand                   brandparam.BrandInfo       `json:"brand"`
	Category                categoryparam.CategoryInfo `json:"category"`
	Gallery                 []entity.ProductGallery    `json:"gallery"`
	IsPublished             bool                       `json:"is_published"`
	IsSecondHand            bool                       `json:"is_second_hand"`
	RequiresFreightShipping bool                       `json:"requires_freight_shipping"`
	Kind                    entity.ProductKind         `json:"kind"`
	SalesChannels           []entity.SalesChannel      `json:"sales_channels"`
	Price                   uint64                     `json:"price,omitempty"`
	Description             string                     `json:"description"`
	Rating                  float64                    `json:"rating"`
	RatingCount             uint                       `json:"rating_count"`
	Discount                []ProductDiscount          `json:"discount"`

	Variants       []StorefrontVariant `json:"variants"`
	Specifications []Specification     `json:"specifications,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`

	SeoTitle       string              `json:"seo_title,omitempty"`
	SeoDescription string              `json:"seo_description,omitempty"`
	FocusKeyword   string              `json:"focus_keyword,omitempty"`
	FAQ            []entity.ProductFAQ `json:"faq,omitempty"`
}

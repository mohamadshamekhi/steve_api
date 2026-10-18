package productparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
)

type GetProductRequest struct {
	ID string `json:"id"`
}
type GetProductResponse struct {
	Product ProductInfo `json:"product"`
}

type ProductInfo struct {
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
	// Price is the fee of a service. Omitted for goods, which price per
	// variant in Variants instead.
	Price       uint64            `json:"price,omitempty"`
	Description string            `json:"description"`
	Rating      float64           `json:"rating"`
	Discount    []ProductDiscount `json:"discount"`

	Variants       []entity.Variant `json:"variants"`
	Specifications []Specification  `json:"specifications,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`

	SeoTitle       string              `json:"seo_title,omitempty"`
	SeoDescription string              `json:"seo_description,omitempty"`
	FocusKeyword   string              `json:"focus_keyword,omitempty"`
	FAQ            []entity.ProductFAQ `json:"faq,omitempty"`
}

type Specification struct {
	Key   string `bson:"key" json:"key"`
	Value string `bson:"value" json:"value"`
	Name  string `bson:"name" json:"name"`
	Unit  string `bson:"unit,omitempty" json:"unit"`
}

type GetProductsByIDsRequest struct {
	IDs []string `json:"ids"`
}

type GetProductsByIDsResponse struct {
	Products map[string]ProductList `json:"products"`
}

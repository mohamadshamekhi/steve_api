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
	ID           string                     `json:"_id"`
	Name         string                     `json:"name"`
	ENTitle      string                     `json:"en_title"`
	Slug         string                     `json:"slug"`
	Brand        brandparam.BrandInfo       `json:"brand"`
	Category     categoryparam.CategoryInfo `json:"category"`
	Gallery      []entity.ProductGallery    `json:"gallery"`
	IsPublished  bool                       `json:"is_published"`
	IsSecondHand bool                       `json:"is_second_hand"`
	Description  string                     `json:"description"`

	Variants       []entity.Variant       `json:"variants"`
	Specifications []entity.Specification `json:"specifications,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

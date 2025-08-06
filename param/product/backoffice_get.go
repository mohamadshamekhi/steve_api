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
	ID           string                     `bson:"_id"`
	Name         string                     `bson:"name"`
	ENTitle      string                     `bson:"en_title"`
	Slug         string                     `bson:"slug"`
	Brand        brandparam.BrandInfo       `bson:"brand"`
	Category     categoryparam.CategoryInfo `bson:"category"`
	Gallery      []entity.ProductGallery    `bson:"gallery"`
	IsPublished  bool                       `bson:"is_published"`
	IsSecondHand bool                       `bson:"is_second_hand"`
	Description  string                     `bson:"description"`

	Variants       []entity.Variant       `bson:"variants"`
	Specifications []entity.Specification `bson:"specifications,omitempty"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at,omitempty"`
}

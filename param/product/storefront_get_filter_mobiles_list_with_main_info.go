package productparam

import (
	"github.com/abolfazlnorzad/steve/entity"
)

type GetFilterMobileListWithMainInfoRequest struct {
	BrandID *string `json:"brand_id"`
	Name    *string `json:"name"`
}

type GetFilterMobileListWithMainInfoResponse struct {
	Products []MobileListWithMainInfo `json:"products"`
}

type MobileListWithMainInfo struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	Slug         string                  `json:"slug"`
	BrandID      string                  `json:"brand_id"`
	CategoryID   string                  `json:"category_id"`
	Gallery      []entity.ProductGallery `json:"gallery"`
	IsPublished  bool                    `json:"is_published"`
	IsSecondHand bool                    `json:"is_second_hand"`
}

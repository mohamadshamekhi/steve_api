package settingparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
)

type CreateBannerRequest struct {
	Title     string     `json:"title"`
	MediaID   string     `json:"media_id"`
	Link      string     `json:"link"`
	Location  string     `json:"location"`
	SortOrder int        `json:"sort_order"`
	IsActive  bool       `json:"is_active"`
	AltText   string     `json:"alt_text"`
	StartsAt  *time.Time `json:"starts_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type CreateBannerResponse struct{}

type UpdateBannerRequest struct {
	ID        string     `json:"id"`
	Title     *string    `json:"title"`
	MediaID   *string    `json:"media_id"`
	Link      *string    `json:"link"`
	Location  *string    `json:"location"`
	SortOrder *int       `json:"sort_order"`
	IsActive  *bool      `json:"is_active"`
	AltText   *string    `json:"alt_text"`
	StartsAt  *time.Time `json:"starts_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}
type UpdateBannerResponse struct{}

type DeleteBannerRequest struct {
	ID string `json:"id"`
}
type DeleteBannerResponse struct{}

type ListBannerRequest struct {
	IsActive *bool `json:"is_active"`
}

type ListBannerResponse struct {
	Banners []entity.Banner `json:"banners"`
}

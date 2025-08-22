package articleparam

import (
	"time"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type AllStorefrontArticleRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}

type AllStorefrontArticleResponse struct {
	Articles   []StorefrontArticle            `json:"articles"`
	Pagination commonparam.PaginationResponse `json:"pagination"`
}
type StorefrontArticle struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	IsPublished bool      `json:"is_published"`
	MediaUrl    string    `json:"media_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type StorefrontArticleDetail struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Body        string    `json:"body"`
	IsPublished bool      `json:"is_published"`
	MediaUrl    string    `json:"media_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

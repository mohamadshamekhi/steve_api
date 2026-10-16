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
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	IsPublished    bool      `json:"is_published"`
	Tags           []string  `json:"tags"`
	BackOfficeID   string    `json:"back_office_id"`
	BackOfficeName string    `json:"back_office_name"`
	MediaUrl       string    `json:"media_url"`
	Alt            string    `json:"alt"`
	ViewCount      int64     `json:"view_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type StorefrontArticleDetail struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Body            string    `json:"body"`
	IsPublished     bool      `json:"is_published"`
	Tags            []string  `json:"tags"`
	MetaDescription string    `json:"meta_description"`
	BackOfficeID    string    `json:"back_office_id"`
	BackOfficeName  string    `json:"back_office_name"`
	MediaUrl        string    `json:"media_url"`
	Alt             string    `json:"alt"`
	ViewCount       int64     `json:"view_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

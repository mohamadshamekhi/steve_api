package searchparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type StorefrontSearchRequest struct {
	Search commonparam.SearchRequest
}

type StorefrontSearchResponse struct {
	Products []StorefrontProductSearchList `json:"products"`
	Articles []StorefrontArticleSearchList `json:"articles"`
}

type StorefrontProductSearchList struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	ENTitle      string                  `json:"en_title"`
	Slug         string                  `json:"slug"`
	Gallery      []entity.ProductGallery `json:"gallery"`
	Rating       float64                 `json:"rating"`
	IsPublished  bool                    `json:"is_published"`
	IsSecondHand bool                    `json:"is_second_hand"`
}

type StorefrontArticleSearchList struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	IsPublished bool      `json:"is_published"`
	MediaUrl    string    `json:"media_url"`
	Alt         string    `json:"alt"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

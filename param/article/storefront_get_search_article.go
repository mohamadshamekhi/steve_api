package articleparam

import (
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	searchparam "github.com/abolfazlnorzad/steve/param/search"
)

type StorefrontGetSearchArticleRequest struct {
	Search commonparam.SearchRequest `json:"search"`
}

type StorefrontGetSearchArticleResponse struct {
	Articles []searchparam.StorefrontArticleSearchList `json:"articles"`
}

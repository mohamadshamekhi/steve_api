package productparam

import (
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	searchparam "github.com/abolfazlnorzad/steve/param/search"
)

type StorefrontGetSearchProductRequest struct {
	Search commonparam.SearchRequest
}

type StorefrontGetSearchProductResponse struct {
	Products []searchparam.StorefrontProductSearchList `json:"products"`
}

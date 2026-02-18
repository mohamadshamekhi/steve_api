package searchservice

import (
	"context"

	articleparam "github.com/abolfazlnorzad/steve/param/article"
	productparam "github.com/abolfazlnorzad/steve/param/product"
)

type ProductProvider interface {
	GetSearchProducts(ctx context.Context, req productparam.StorefrontGetSearchProductRequest) (productparam.StorefrontGetSearchProductResponse, error)
}

type ArticleProvider interface {
	GetSearchArticles(ctx context.Context, req articleparam.StorefrontGetSearchArticleRequest) (articleparam.StorefrontGetSearchArticleResponse, error)
}

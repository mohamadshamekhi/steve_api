package searchservice

import (
	"context"

	articleparam "github.com/abolfazlnorzad/steve/param/article"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	searchparam "github.com/abolfazlnorzad/steve/param/search"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StorefrontService struct {
	productProvider ProductProvider
	articleProvider ArticleProvider
}

func NewStorefrontService(productProvider ProductProvider, articleProvider ArticleProvider) StorefrontService {
	return StorefrontService{
		productProvider: productProvider,
		articleProvider: articleProvider,
	}
}

func (s StorefrontService) SearchInProductAndBlog(ctx context.Context, req searchparam.StorefrontSearchRequest) (searchparam.StorefrontSearchResponse, error) {
	const op = "searchservice.SearchInProductAndBlog"

	products, err := s.productProvider.GetSearchProducts(ctx, productparam.StorefrontGetSearchProductRequest{Search: req.Search})
	if err != nil {
		return searchparam.StorefrontSearchResponse{}, richerror.New(op).WithWarpError(err)
	}

	articles, aErr := s.articleProvider.GetSearchArticles(ctx, articleparam.StorefrontGetSearchArticleRequest{Search: req.Search})
	if aErr != nil {
		return searchparam.StorefrontSearchResponse{}, richerror.New(op).WithWarpError(aErr)
	}

	var listProduct []searchparam.StorefrontProductSearchList
	var listArticle []searchparam.StorefrontArticleSearchList

	for _, product := range products.Products {
		listProduct = append(listProduct, searchparam.StorefrontProductSearchList{
			ID:           product.ID,
			Name:         product.Name,
			ENTitle:      product.ENTitle,
			Slug:         product.Slug,
			Gallery:      product.Gallery,
			Rating:       product.Rating,
			IsPublished:  product.IsPublished,
			IsSecondHand: product.IsSecondHand,
		})
	}

	for _, article := range articles.Articles {
		listArticle = append(listArticle, searchparam.StorefrontArticleSearchList{
			ID:          article.ID,
			Title:       article.Title,
			Slug:        article.Slug,
			IsPublished: article.IsPublished,
			MediaUrl:    article.MediaUrl,
			Alt:         article.Alt,
			CreatedAt:   article.CreatedAt,
			UpdatedAt:   article.UpdatedAt,
		})
	}

	return searchparam.StorefrontSearchResponse{
		Products: listProduct,
		Articles: listArticle,
	}, nil
}

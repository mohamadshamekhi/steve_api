package productservice

import (
	"context"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	searchparam "github.com/abolfazlnorzad/steve/param/search"
)

func (s StorefrontService) GetSearchProducts(ctx context.Context, req productparam.StorefrontGetSearchProductRequest) (productparam.StorefrontGetSearchProductResponse, error) {
	const op = "productservice.GetSearchProduct"

	resp, err := s.repo.GetSearchProduct(ctx, req.Search)
	if err != nil {
		return productparam.StorefrontGetSearchProductResponse{}, nil
	}

	var list []searchparam.StorefrontProductSearchList

	for _, product := range resp.Products {
		list = append(list, searchparam.StorefrontProductSearchList{
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

	return productparam.StorefrontGetSearchProductResponse{
		Products: list,
	}, nil
}

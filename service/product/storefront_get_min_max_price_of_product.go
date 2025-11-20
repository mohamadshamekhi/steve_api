package productservice

import (
	"context"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s StorefrontService) GetMinAdnMaxPriceOfProductWithSpecialCategory(ctx context.Context, req productparam.GetMinAdnMaxPriceOfProductWithSpecialCategoryRequest) (productparam.GetMinAdnMaxPriceOfProductWithSpecialCategoryResponse, error) {
	const op = "productservice.GetMinAdnMaxPriceOfProductWithSpecialCategory"

	minPrice, maxPrice, err := s.repo.GetMinAdnMaxPriceOfProductWithSpecialCategory(ctx, req.CategoryID)
	if err != nil {
		return productparam.GetMinAdnMaxPriceOfProductWithSpecialCategoryResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.GetMinAdnMaxPriceOfProductWithSpecialCategoryResponse{
		MaxPrice: maxPrice,
		MinPrice: minPrice,
	}, nil
}

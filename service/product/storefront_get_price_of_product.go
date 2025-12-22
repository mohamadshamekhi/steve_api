package productservice

import (
	"context"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s BackofficeService) GetPricesOfProducts(ctx context.Context, req productparam.GetPriceOfProductsRequest) (productparam.GetPriceOfProductsResponse, error) {
	const op = "productservice.GetPriceOfProduct"
	p, err := s.repo.GetPricesOfProducts(ctx, req.ProductIDs, req.VariantIDs)

	if err != nil {
		return productparam.GetPriceOfProductsResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.GetPriceOfProductsResponse{
		Price: p,
	}, nil
}

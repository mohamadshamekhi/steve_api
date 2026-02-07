package productservice

import (
	"context"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s StorefrontService) IncrementSalesCount(ctx context.Context, req productparam.UpdateProductSalesCountRequest) (productparam.UpdateProductSalesCountResponse, error) {
	const op = "productservice.IncrementSalesCount"
	if err := s.repo.UpdateProductSalesCount(ctx, req.ProductID, req.Quantity); err != nil {
		return productparam.UpdateProductSalesCountResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.UpdateProductSalesCountResponse{}, nil
}

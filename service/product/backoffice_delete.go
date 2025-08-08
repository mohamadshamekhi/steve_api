package productservice

import (
	"context"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s BackofficeService) DeleteProduct(ctx context.Context, req productparam.DeleteProductRequest) (productparam.DeleteProductResponse, error) {
	const op = "productservice.DeleteProduct"
	p, gErr := s.repo.FindByID(ctx, req.ID)
	if gErr != nil {
		return productparam.DeleteProductResponse{}, richerror.New(op).WithWarpError(gErr)
	}
	if err := s.repo.Delete(ctx, p); err != nil {
		return productparam.DeleteProductResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.DeleteProductResponse{}, nil
}

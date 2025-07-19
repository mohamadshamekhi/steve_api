package productservice

import (
	"context"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s BackofficeService) GetProduct(ctx context.Context, req productparam.GetProductRequest) (productparam.GetProductResponse, error) {
	const op = "productservice.GetProduct"

	p, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return productparam.GetProductResponse{}, richerror.New(op).WithWarpError(err)
	}

	return productparam.GetProductResponse{Product: p}, nil
}

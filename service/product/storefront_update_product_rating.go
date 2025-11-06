package productservice

import (
	"context"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s StorefrontService) UpdateProductRating(ctx context.Context, req productparam.UpdateProductRatingRequest) (productparam.UpdateProductRatingResponse, error) {
	const op = "productservice.UpdateProductRating"

	if err := s.repo.UpdateProductAvrRating(ctx, req.ProductID, req.AverageRating); err != nil {
		return productparam.UpdateProductRatingResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.UpdateProductRatingResponse{}, nil
}

package productservice

import (
	"context"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s StorefrontService) GetFilterableAttributes(ctx context.Context, req productparam.GetFilterForCategoryRequest) (productparam.GetFilterForCategoryResponse, error) {
	const op = "product.GetFilterableAttributes"

	ag, err := s.repo.GetFilterableAttributes(ctx, req.CategoryID)
	if err != nil {
		return productparam.GetFilterForCategoryResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.GetFilterForCategoryResponse{
		AggregatedFilterOption: ag,
	}, nil
}

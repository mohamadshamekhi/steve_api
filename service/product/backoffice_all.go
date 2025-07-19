package productservice

import (
	"context"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s BackofficeService) GetAllProducts(ctx context.Context, req productparam.GetAllProductsRequest) (productparam.GetAllProductsResponse, error) {
	const op = "productservice.GetAllProducts"

	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return productparam.GetAllProductsResponse{}, richerror.New(op).WithWarpError(err)
	}

	return productparam.GetAllProductsResponse{
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)},
		Products: resp.Products,
	}, nil
}

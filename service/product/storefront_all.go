package productservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s StorefrontService) GetFilteredProduct(ctx context.Context, req productparam.StorefrontGetAllProductRequest) (productparam.StorefrontGetAllProductResponse, error) {
	const op = "productservice.GetFilteredProduct"

	resp, err := s.repo.GetFilteredProducts(ctx, &req.Filter, req.Pagination)
	if err != nil {
		return productparam.StorefrontGetAllProductResponse{}, richerror.New(op).WithWarpError(err)
	}
	res, err := s.discountProvider.GetProductIDsThatHaveActiveAutoDiscount(ctx, discountparam.GetProductIDsThatHaveActiveAutoDiscountRequest{})
	if err != nil {
		return productparam.StorefrontGetAllProductResponse{}, richerror.New(op).WithWarpError(err)
	}

	list := []productparam.StorefrontProductList{}

	for _, product := range resp.Products {
		discount := entity.Discount{}
		if d, ok := res.DiscountMap[product.ID]; ok {
			discount = d
		}
		list = append(list, productparam.StorefrontProductList{
			ID:           product.ID,
			Name:         product.Name,
			ENTitle:      product.ENTitle,
			Slug:         product.Slug,
			Gallery:      product.Gallery,
			IsPublished:  product.IsPublished,
			IsSecondHand: product.IsSecondHand,
			Discount:     discount,
			Variant:      product.Variants,
			Rating:       product.Rating,
		})
	}

	return productparam.StorefrontGetAllProductResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(resp.TotalCount),
		},
		Products: list,
	}, nil
}

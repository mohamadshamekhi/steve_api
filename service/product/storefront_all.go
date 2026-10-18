package productservice

import (
	"context"

	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"golang.org/x/sync/errgroup"
)

func (s StorefrontService) GetFilteredProduct(ctx context.Context, req productparam.StorefrontGetAllProductRequest) (productparam.StorefrontGetAllProductResponse, error) {
	const op = "productservice.GetFilteredProduct"

	if req.Filter.CategorySlug != nil {
		category, cErr := s.categoryProvider.GetCategoryBySlug(ctx, categoryparam.GetCategoryBySlugRequest{Slug: *req.Filter.CategorySlug})
		if cErr != nil {
			logger.L().Error("failed to get category by slug", "category_slug", *req.Filter.CategorySlug, "err", cErr)
		}
		req.Filter.CategoryID = &category.Category.ID
	}

	resp, err := s.repo.GetFilteredProducts(ctx, &req.Filter, req.Pagination)
	if err != nil {
		return productparam.StorefrontGetAllProductResponse{}, richerror.New(op).WithWarpError(err)
	}
	res, err := s.discountProvider.GetProductIDsThatHaveActiveAutoDiscount(ctx, discountparam.GetProductIDsThatHaveActiveAutoDiscountRequest{})
	if err != nil {
		return productparam.StorefrontGetAllProductResponse{}, richerror.New(op).WithWarpError(err)
	}

	list := make([]productparam.StorefrontProductList, len(resp.Products))

	g, gCtx := errgroup.WithContext(ctx)

	for i, p := range resp.Products {
		index := i
		product := p

		g.Go(func() error {

			var catRes categoryparam.GetCategoryResponse
			var cErr error

			catRes, cErr = s.categoryProvider.GetCategory(gCtx, categoryparam.GetCategoryRequest{ID: product.CategoryID})
			if cErr != nil {
				logger.L().Error("failed to get category", "product_id", product.ID, "err", cErr)
			}

			productSpecMap := make(map[string]string)
			for _, s := range product.Specifications {
				productSpecMap[s.Key] = s.Value
			}
			var specs []productparam.Specification
			if catRes.Category.ID != "" {
				for _, attr := range catRes.Category.AttributeSet {
					if val, ok := productSpecMap[attr.Key]; ok {
						specs = append(specs, productparam.Specification{
							Key:   attr.Key,
							Value: val,
							Name:  attr.Name,
							Unit:  attr.Unit,
						})
					}
				}
			}
			// we don't need mutex because every goroutine write on spsefic index
			list[index] = productparam.StorefrontProductList{
				ID:             product.ID,
				Name:           product.Name,
				ENTitle:        product.ENTitle,
				Slug:           product.Slug,
				Gallery:        product.Gallery,
				IsPublished:    product.IsPublished,
				IsSecondHand:   product.IsSecondHand,
				Discount:       res.DiscountMap[product.ID],
				Variant:        productparam.ToStorefrontVariants(product.Variants),
				Rating:         product.Rating,
				Specifications: specs,
			}
			return nil
		})
	}
	g.Wait()

	return productparam.StorefrontGetAllProductResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(resp.TotalCount),
		},
		Products: list,
	}, nil
}

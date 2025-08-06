package productservice

import (
	"context"
	"fmt"

	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s BackofficeService) GetAllProducts(ctx context.Context, req productparam.GetAllProductsRequest) (productparam.GetAllProductsResponse, error) {
	const op = "productservice.GetAllProducts"

	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return productparam.GetAllProductsResponse{}, richerror.New(op).WithWarpError(err)
	}
	var productList []productparam.ProductList

	var CateIds []string
	for _, p := range resp.Products {
		CateIds = append(CateIds, p.CategoryID)
	}
	var brandIDs []string
	for _, p := range resp.Products {
		brandIDs = append(brandIDs, p.BrandID)
	}
	brandsRes, bErr := s.brandProvider.GetBrandsByIDs(ctx, brandparam.GetBrandsByIDsRequest{IDs: brandIDs})
	if bErr != nil {
		logger.L().Error(fmt.Sprintf("cann't get brand's of product-err%s - %v", op, bErr))
	}
	cateRes, cErr := s.categoryProvider.GetCategoriesByIDs(ctx, categoryparam.GetCategoriesByIDsRequest{IDs: CateIds})
	if bErr != nil {
		logger.L().Error(fmt.Sprintf("cann't get category's of product-err%s - %v", op, cErr))
	}
	for _, p := range resp.Products {

		productList = append(productList, productparam.ProductList{
			ID:           p.ID,
			Name:         p.Name,
			EnTitle:      p.ENTitle,
			Brand:        brandsRes.Brands[p.BrandID].Name,
			Category:     cateRes.Categories[p.CategoryID].Name,
			IsPublished:  p.IsPublished,
			IsSecondHand: p.IsSecondHand,
			Slug:         p.Slug,
		})
	}

	return productparam.GetAllProductsResponse{
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)},
		Products: productList,
	}, nil
}

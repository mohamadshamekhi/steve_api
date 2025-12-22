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

func (s BackofficeService) GetProduct(ctx context.Context, req productparam.GetProductRequest) (productparam.GetProductResponse, error) {
	const op = "productservice.GetProduct"

	p, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return productparam.GetProductResponse{}, richerror.New(op).WithWarpError(err)
	}

	productInfo := productparam.ProductInfo{
		ID:             p.ID,
		Name:           p.Name,
		ENTitle:        p.ENTitle,
		Slug:           p.Slug,
		Brand:          brandparam.BrandInfo{},
		Category:       categoryparam.CategoryInfo{},
		Gallery:        p.Gallery,
		IsPublished:    p.IsPublished,
		IsSecondHand:   p.IsSecondHand,
		Description:    p.Description,
		Variants:       p.Variants,
		Specifications: p.Specifications,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
	cr, cErr := s.categoryProvider.GetCategory(ctx, categoryparam.GetCategoryRequest{ID: p.CategoryID})
	if cErr != nil {
		logger.L().Error(fmt.Sprintf("failed to get category: product_id : %s - err : %v", p.ID, cErr))
	}
	productInfo.Category = categoryparam.CategoryInfo{
		ID:           cr.Category.ID,
		Name:         cr.Category.Name,
		Description:  cr.Category.Description,
		URL:          cr.Category.URL,
		Level:        cr.Category.Level,
		Slug:         cr.Category.Slug,
		Ancestors:    cr.Category.Ancestors,
		AttributeSet: cr.Category.AttributeSet,
	}

	bres, bErr := s.brandProvider.GetBrand(ctx, brandparam.GetBrandRequest{ID: p.BrandID})
	if bErr != nil {
		logger.L().Error(fmt.Sprintf("failed to get banner: product_id : %s - err : %v", p.ID, bErr))
	}

	productInfo.Brand = brandparam.BrandInfo{
		ID:          bres.ID,
		Name:        bres.Name,
		Description: bres.Description,
		URL:         bres.Url,
	}
	return productparam.GetProductResponse{Product: productInfo}, nil
}

func (s BackofficeService) GetProductCategoriesForCartItems(ctx context.Context, req productparam.GetProductCategoriesForCartItemsRequest) (productparam.GetProductCategoriesForCartItemsResponse, error) {
	const op = "productservice.GetProductCategoriesForCartItems"
	res, err := s.repo.GetProductCategoriesForCartItems(ctx, req.Items)
	if err != nil {
		return productparam.GetProductCategoriesForCartItemsResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.GetProductCategoriesForCartItemsResponse{
		CategoryIDs: res,
	}, nil
}

func (s BackofficeService) GetProductsByIDs(ctx context.Context, req productparam.GetProductsByIDsRequest) (productparam.GetProductsByIDsResponse, error) {
	const op = "productservice.GetProductsByIDs"
	products, err := s.repo.GetProductsByIDs(ctx, req.IDs)
	if err != nil {
		return productparam.GetProductsByIDsResponse{}, richerror.New(op).WithWarpError(err)
	}
	res := make(map[string]productparam.ProductList, len(products))

	for _, p := range products {
		res[p.ID] = productparam.ProductList{
			ID:           p.ID,
			Name:         p.Name,
			EnTitle:      p.ENTitle,
			Brand:        p.BrandID,
			Category:     p.CategoryID,
			IsPublished:  p.IsPublished,
			IsSecondHand: p.IsSecondHand,
			Slug:         p.Slug,
			ThumbnailURL: "",
		}
	}
	return productparam.GetProductsByIDsResponse{Products: res}, nil
}

func (s BackofficeService) FindProductWithSpecialCategory(ctx context.Context, req productparam.GetProductWithSpecialCategoryRequest) (productparam.GetProductWithSpecialCategoryResponse, error) {
	const op = "productservice.FindProductWithSpecialCategory"

	pid, err := s.repo.FindProductsIDsByCategory(ctx, req.CategoryID)
	if err != nil {
		return productparam.GetProductWithSpecialCategoryResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.GetProductWithSpecialCategoryResponse{
		ProductID: pid,
	}, nil
}

func (s BackofficeService) GetProductIDsByCategories(ctx context.Context, req productparam.GetProductIDsByCategoriesRequest) (productparam.GetProductIDsByCategoriesResponse, error) {
	const op = "productservice.GetProductIDsByCategories"
	resp, err := s.repo.GetProductIDsByCategories(ctx, req.CatIDs)
	if err != nil {
		return productparam.GetProductIDsByCategoriesResponse{}, err
	}
	return productparam.GetProductIDsByCategoriesResponse{ProdsByCat: resp}, nil
}

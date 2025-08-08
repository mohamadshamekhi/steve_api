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

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

// ExportProducts loads every non-deleted product (no pagination) and shapes it into
// export rows, resolving brand/category names the same way GetAllProducts does.
func (s BackofficeService) ExportProducts(ctx context.Context, req productparam.ExportProductsRequest) (productparam.ExportProductsResponse, error) {
	const op = "productservice.ExportProducts"

	products, err := s.repo.AllForExport(ctx, req.Filter)
	if err != nil {
		return productparam.ExportProductsResponse{}, richerror.New(op).WithWarpError(err)
	}

	var brandIDs, categoryIDs []string
	for _, p := range products {
		brandIDs = append(brandIDs, p.BrandID)
		categoryIDs = append(categoryIDs, p.CategoryID)
	}

	brandsRes, bErr := s.brandProvider.GetBrandsByIDs(ctx, brandparam.GetBrandsByIDsRequest{IDs: brandIDs})
	if bErr != nil {
		logger.L().Error(fmt.Sprintf("cann't get brand's of product-err%s - %v", op, bErr))
	}
	cateRes, cErr := s.categoryProvider.GetCategoriesByIDs(ctx, categoryparam.GetCategoriesByIDsRequest{IDs: categoryIDs})
	if cErr != nil {
		logger.L().Error(fmt.Sprintf("cann't get category's of product-err%s - %v", op, cErr))
	}

	rows := make([]productparam.ProductExportRow, 0, len(products))
	for _, p := range products {
		var minPrice, maxPrice, totalQty uint64
		for i, v := range p.Variants {
			if i == 0 || v.Price < minPrice {
				minPrice = v.Price
			}
			if v.Price > maxPrice {
				maxPrice = v.Price
			}
			totalQty += v.Quantity
		}

		rows = append(rows, productparam.ProductExportRow{
			ID:            p.ID,
			Name:          p.Name,
			EnTitle:       p.ENTitle,
			Slug:          p.Slug,
			Brand:         brandsRes.Brands[p.BrandID].Name,
			Category:      cateRes.Categories[p.CategoryID].Name,
			IsPublished:   p.IsPublished,
			IsSecondHand:  p.IsSecondHand,
			MinPrice:      minPrice,
			MaxPrice:      maxPrice,
			TotalQuantity: totalQty,
			VariantCount:  len(p.Variants),
			SalesCount:    p.SalesCount,
			Rating:        p.Rating,
			CreatedAt:     p.CreatedAt,
		})
	}

	return productparam.ExportProductsResponse{Rows: rows}, nil
}

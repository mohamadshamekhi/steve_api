package productservice

import (
	"context"

	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"golang.org/x/sync/errgroup"
)

func (s BackofficeService) GetProduct(ctx context.Context, req productparam.GetProductRequest) (productparam.GetProductResponse, error) {
	const op = "productservice.GetProduct"

	p, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return productparam.GetProductResponse{}, richerror.New(op).WithWarpError(err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	var (
		catRes   categoryparam.GetCategoryResponse
		brandRes brandparam.GetBrandResponse
	)

	g.Go(func() error {
		var cErr error
		catRes, cErr = s.categoryProvider.GetCategory(gCtx, categoryparam.GetCategoryRequest{ID: p.CategoryID})
		if cErr != nil {
			logger.L().Error("failed to get category", "product_id", p.ID, "err", cErr)
		}
		return nil
	})

	g.Go(func() error {
		var bErr error
		brandRes, bErr = s.brandProvider.GetBrand(gCtx, brandparam.GetBrandRequest{ID: p.BrandID})
		if bErr != nil {
			logger.L().Error("failed to get brand", "product_id", p.ID, "err", bErr)
		}
		return nil
	})

	_ = g.Wait()

	productSpecMap := make(map[string]string)
	for _, s := range p.Specifications {
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

	productInfo := productparam.ProductInfo{
		ID:                      p.ID,
		Name:                    p.Name,
		ENTitle:                 p.ENTitle,
		Slug:                    p.Slug,
		Gallery:                 p.Gallery,
		IsPublished:             p.IsPublished,
		IsSecondHand:            p.IsSecondHand,
		RequiresFreightShipping: p.RequiresFreightShipping,
		Description:             p.Description,
		Kind:                    p.GetKind(),
		SalesChannels:           p.GetSalesChannels(),
		Price:                   p.ServicePrice(),
		Variants:                p.Variants,
		CreatedAt:               p.CreatedAt,
		UpdatedAt:               p.UpdatedAt,
		Specifications:          specs,
		SeoTitle:                p.SeoTitle,
		SeoDescription:          p.SeoDescription,
		FocusKeyword:            p.FocusKeyword,
		FAQ:                     p.FAQ,
	}

	if catRes.Category.ID != "" {
		productInfo.Category = categoryparam.CategoryInfo{
			ID:           catRes.Category.ID,
			Name:         catRes.Category.Name,
			Description:  catRes.Category.Description,
			URL:          catRes.Category.URL,
			Level:        catRes.Category.Level,
			Slug:         catRes.Category.Slug,
			Ancestors:    catRes.Category.Ancestors,
			AttributeSet: catRes.Category.AttributeSet,
		}
	}

	if brandRes.ID != "" {
		productInfo.Brand = brandparam.BrandInfo{
			ID:          brandRes.ID,
			Name:        brandRes.Name,
			Slug:        brandRes.Slug,
			Description: brandRes.Description,
			URL:         brandRes.Url,
		}
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
			ID:                      p.ID,
			Name:                    p.Name,
			EnTitle:                 p.ENTitle,
			Brand:                   p.BrandID,
			Category:                p.CategoryID,
			IsPublished:             p.IsPublished,
			IsSecondHand:            p.IsSecondHand,
			RequiresFreightShipping: p.RequiresFreightShipping,
			Slug:                    p.Slug,
			ThumbnailURL:            "",
			Kind:                    p.GetKind(),
			SalesChannels:           p.GetSalesChannels(),
			Price:                   p.ServicePrice(),
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

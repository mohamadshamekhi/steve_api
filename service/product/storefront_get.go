package productservice

import (
	"context"
	"errors"

	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"golang.org/x/sync/errgroup"
)

func (s StorefrontService) GetProduct(ctx context.Context, req productparam.StorefrontGetProductRequest) (productparam.StorefrontGetProductResponse, error) {
	const op = "StorefrontService.GetProduct"

	p, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return productparam.StorefrontGetProductResponse{}, richerror.New(op).WithWarpError(err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	var (
		catRes   categoryparam.GetCategoryResponse
		brandRes brandparam.GetBrandResponse
		disRes   discountparam.GetProductIDsThatHaveActiveAutoDiscountResponse
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
	g.Go(func() error {
		disRes, err = s.discountProvider.GetProductIDsThatHaveActiveAutoDiscount(ctx, discountparam.GetProductIDsThatHaveActiveAutoDiscountRequest{})
		if err != nil {
			logger.L().Error("failed to get product ids", "product_id", p.ID, "err", err)
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

	res := productparam.StorefrontProductInfo{
		ID:             p.ID,
		Name:           p.Name,
		ENTitle:        p.ENTitle,
		Slug:           p.Slug,
		Gallery:        p.Gallery,
		IsPublished:    p.IsPublished,
		IsSecondHand:   p.IsSecondHand,
		Description:    p.Description,
		Variants:       productparam.ToStorefrontVariants(p.Variants),
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
		Discount:       disRes.DiscountMap[p.ID],
		Specifications: specs,
		SeoTitle:       p.SeoTitle,
		SeoDescription: p.SeoDescription,
		FocusKeyword:   p.FocusKeyword,
		FAQ:            p.FAQ,
	}

	if catRes.Category.ID != "" {
		res.Category = categoryparam.CategoryInfo{
			ID:           catRes.Category.ID,
			Name:         catRes.Category.Name,
			Slug:         catRes.Category.Slug,
			URL:          catRes.Category.URL,
			AttributeSet: catRes.Category.AttributeSet,
		}
	}

	if brandRes.ID != "" {
		res.Brand = brandparam.BrandInfo{
			ID:          brandRes.ID,
			Name:        brandRes.Name,
			Slug:        brandRes.Slug,
			Description: brandRes.Description,
		}
	}

	return productparam.StorefrontGetProductResponse{Product: res}, nil
}
func (s StorefrontService) GetProductBySlug(ctx context.Context, req productparam.StorefrontGetProductBySlugRequest) (productparam.StorefrontGetProductBySlugResponse, error) {
	const op = "StorefrontService.GetProductBySlug"

	p, err := s.repo.FindBySlug(ctx, req.Slug)
	if err != nil {
		// A renamed product is still served from its old slug; the response
		// carries the current slug, which the storefront 301-redirects to.
		var rErr *richerror.RichError
		if !errors.As(err, &rErr) || rErr.Kind() != richerror.KindNotFound {
			return productparam.StorefrontGetProductBySlugResponse{}, richerror.New(op).WithWarpError(err)
		}
		prev, pErr := s.repo.FindByPreviousSlug(ctx, req.Slug)
		if pErr != nil {
			return productparam.StorefrontGetProductBySlugResponse{}, richerror.New(op).WithWarpError(err)
		}
		p = prev
	}

	g, gCtx := errgroup.WithContext(ctx)

	var (
		catRes   categoryparam.GetCategoryResponse
		brandRes brandparam.GetBrandResponse
		disRes   discountparam.GetProductIDsThatHaveActiveAutoDiscountResponse
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
	g.Go(func() error {
		disRes, err = s.discountProvider.GetProductIDsThatHaveActiveAutoDiscount(ctx, discountparam.GetProductIDsThatHaveActiveAutoDiscountRequest{})
		if err != nil {
			logger.L().Error("failed to get product ids", "product_id", p.ID, "err", err)
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

	res := productparam.StorefrontProductInfo{
		ID:             p.ID,
		Name:           p.Name,
		ENTitle:        p.ENTitle,
		Slug:           p.Slug,
		Gallery:        p.Gallery,
		IsPublished:    p.IsPublished,
		IsSecondHand:   p.IsSecondHand,
		Description:    p.Description,
		Variants:       productparam.ToStorefrontVariants(p.Variants),
		Discount:       disRes.DiscountMap[p.ID],
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
		Specifications: specs,
		SeoTitle:       p.SeoTitle,
		SeoDescription: p.SeoDescription,
		FocusKeyword:   p.FocusKeyword,
		FAQ:            p.FAQ,
	}

	if catRes.Category.ID != "" {
		res.Category = categoryparam.CategoryInfo{
			ID:           catRes.Category.ID,
			Name:         catRes.Category.Name,
			Slug:         catRes.Category.Slug,
			URL:          catRes.Category.URL,
			AttributeSet: catRes.Category.AttributeSet,
		}
	}

	if brandRes.ID != "" {
		res.Brand = brandparam.BrandInfo{
			ID:          brandRes.ID,
			Name:        brandRes.Name,
			Slug:        brandRes.Slug,
			Description: brandRes.Description,
		}
	}

	return productparam.StorefrontGetProductBySlugResponse{Product: res}, nil
}

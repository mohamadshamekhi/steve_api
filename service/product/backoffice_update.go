package productservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/abolfazlnorzad/steve/pkg/slug"
)

func (s BackofficeService) UpdateProduct(ctx context.Context, req productparam.UpdateProductRequest) (productparam.UpdateProductResponse, error) {
	const op = "productservice.UpdateProduct"

	p, gErr := s.repo.FindByID(ctx, req.ID)
	if gErr != nil {
		return productparam.UpdateProductResponse{}, richerror.New(op).WithWarpError(gErr)
	}
	if req.ENTitle != nil && *req.ENTitle != p.ENTitle {
		p.ENTitle = *req.ENTitle

		slug := slug.GenerateSlug(p.ENTitle)
		counter := 1
		for {
			existingArticle, err := s.repo.FindBySlug(ctx, slug)
			if err != nil {
				var rErr *richerror.RichError
				if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
					break
				}
				return productparam.UpdateProductResponse{}, richerror.New(op).WithWarpError(err)
			}

			if existingArticle.ID == "" {
				break
			}

			slug = fmt.Sprintf("%s-%d", slug, counter)
			counter++
		}

		p.Slug = slug
	}

	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.BrandID != nil {
		p.BrandID = *req.BrandID
	}

	if req.CategoryID != nil {
		p.CategoryID = *req.CategoryID
	}
	if req.IsSecondHand != nil {
		p.IsSecondHand = *req.IsSecondHand
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return productparam.UpdateProductResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.UpdateProductResponse{Product: p}, nil
}

func (s BackofficeService) UpdateProductVariantQuantityAfterRegisterOrder(ctx context.Context, req productparam.UpdateProductVariantQuantityAfterRegisterOrderRequest) (productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse, error) {
	const op = "productservice.UpdateProductVariantQuantityAfterRegisterOrder"

	for _, item := range req.Order.Items {
		p, _ := s.repo.FindByID(ctx, item.ProductID)
		var v entity.Variant
		for _, variant := range p.Variants {
			if variant.ID == item.VariantID {
				v = variant
			}
		}
		if uint(v.Quantity) < item.Quantity {
			return productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse{}, richerror.New(op).WithKind(richerror.KindUnexpected)
		}
		finalQuantity := uint(v.Quantity) - item.Quantity
		if err := s.repo.UpdateQuantity(ctx, p.ID, v, finalQuantity); err != nil {
			return productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse{}, richerror.New(op).WithWarpError(err)
		}
	}
	return productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse{}, nil
}

func (s BackofficeService) PublishProduct(ctx context.Context, req productparam.PublishProductRequest) (productparam.PublishProductResponse, error) {
	const op = "productservice.PublishProduct"
	product, err := s.repo.FindByID(ctx, req.ProductID)
	if err != nil {
		return productparam.PublishProductResponse{}, richerror.New(op).WithWarpError(err)
	}
	if product.IsPublished {
		return productparam.PublishProductResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("product.publish.already_published")
	}
	if err = s.runPrePublishChecks(ctx, &product); err != nil {
		return productparam.PublishProductResponse{}, richerror.New(op).WithWarpError(err)
	}
	product.IsPublished = true
	if err := s.repo.Update(ctx, product); err != nil {
		return productparam.PublishProductResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.PublishProductResponse{}, nil
}

func (s BackofficeService) runPrePublishChecks(ctx context.Context, product *entity.Product) error {
	const op = "product.service.runPrePublishChecks"

	if len(product.Variants) == 0 {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("product.publish.no_variants")
	}

	if len(product.Gallery) == 0 {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("product.publish.no_gallery_images")
	}

	hasPrimaryImage := false
	for _, galleryItem := range product.Gallery {
		if galleryItem.IsPrimary {
			hasPrimaryImage = true
			break
		}
	}
	if !hasPrimaryImage {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("product.publish.no_primary_image")
	}

	if product.CategoryID == "" {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("product.publish.no_category_assigned")
	}

	categoryResp, err := s.categoryProvider.GetCategory(ctx, categoryparam.GetCategoryRequest{ID: product.CategoryID})
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("product.publish.category_not_found")
	}

	if categoryResp.Category.AttributeSet == nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("product.publish.attribute_set_not_found")
	}

	providedSpecs := make(map[string]struct{})
	for _, spec := range product.Specifications {
		providedSpecs[spec.Key] = struct{}{}
	}

	var missingSpecs []string
	for _, requiredDef := range categoryResp.Category.AttributeSet {
		if _, ok := providedSpecs[requiredDef.Key]; !ok {
			missingSpecs = append(missingSpecs, requiredDef.Name)
		}
	}

	if len(missingSpecs) > 0 {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("product.publish.missing_specifications").
			WithMeta(map[string]any{"missing_specs": missingSpecs})
	}

	return nil
}

func (s BackofficeService) UnPublishProduct(ctx context.Context, req productparam.UnPublishProductRequest) (productparam.UnPublishProductResponse, error) {
	const op = "productservice.UnPublishProduct"
	product, err := s.repo.FindByID(ctx, req.ProductID)
	if err != nil {
		return productparam.UnPublishProductResponse{}, richerror.New(op).WithWarpError(err)
	}
	if !product.IsPublished {
		return productparam.UnPublishProductResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("product.publish.already_unpublished")
	}

	product.IsPublished = false
	if err := s.repo.Update(ctx, product); err != nil {
		return productparam.UnPublishProductResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.UnPublishProductResponse{}, nil
}

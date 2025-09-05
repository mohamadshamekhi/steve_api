package productservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
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

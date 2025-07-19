package productservice

import (
	"context"
	"errors"
	"fmt"
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

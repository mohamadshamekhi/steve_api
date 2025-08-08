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

func (s BackofficeService) Create(ctx context.Context, req productparam.CreateProductRequest) (productparam.CreateProductResponse, error) {
	const op = "productservice.Create"

	p := entity.Product{
		Name:         req.Name,
		ENTitle:      req.ENTitle,
		BrandID:      req.BrandID,
		CategoryID:   req.CategoryID,
		IsSecondHand: *req.IsSecondHand,
		Description:  req.Description,
	}
	slg := slug.GenerateSlug(p.ENTitle)
	counter := 1
	for {
		existingArticle, err := s.repo.FindBySlug(ctx, slg)
		if err != nil {
			var rErr *richerror.RichError
			if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
				break
			}
			return productparam.CreateProductResponse{}, richerror.New(op).WithWarpError(err)
		}

		if existingArticle.ID == "" {
			break
		}

		slg = fmt.Sprintf("%s-%d", slg, counter)
		counter++
	}
	p.Slug = slg
	p.BeforeCreate()

	if err := s.repo.Create(ctx, p); err != nil {
		return productparam.CreateProductResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.CreateProductResponse{Product: p}, nil
}

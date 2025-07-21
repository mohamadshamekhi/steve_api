package productservice

import (
	"context"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s BackofficeService) UpdateSpecifications(ctx context.Context, req productparam.AddProductAttributeRequest) (productparam.AddProductAttributeResponse, error) {
	const op = "product.service.UpdateSpecifications"

	product, err := s.repo.FindByID(ctx, req.ProductID)
	if err != nil {
		return productparam.AddProductAttributeResponse{}, richerror.New(op).WithWarpError(err)
	}

	resp, err := s.categoryProvider.GetCategory(ctx, categoryparam.GetCategoryRequest{
		ID: product.CategoryID,
	})
	if err != nil {
		return productparam.AddProductAttributeResponse{}, richerror.New(op).WithWarpError(err)
	}

	validKeys := make(map[string]struct{})
	for _, def := range resp.Category.AttributeSet {
		validKeys[def.Key] = struct{}{}
	}

	for _, spec := range req.Attrs {
		if _, ok := validKeys[spec.Key]; !ok {
			return productparam.AddProductAttributeResponse{}, richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("product.invalid_specification_key").
				WithMeta(map[string]any{"invalid_key": spec.Key})
		}
	}

	if err := s.repo.UpdateSpecifications(ctx, req.ProductID, req.Attrs); err != nil {
		return productparam.AddProductAttributeResponse{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	product.Specifications = req.Attrs
	return productparam.AddProductAttributeResponse{}, nil
}

func (s BackofficeService) UpdateSingleAttribute(ctx context.Context, req productparam.UpdateSingleAttributeRequest) (productparam.UpdateSingleAttributeResponse, error) {
	const op = "productservice.UpdateSingleAttribute"
	product, err := s.repo.FindByID(ctx, req.ProductID)
	if err != nil {
		return productparam.UpdateSingleAttributeResponse{}, richerror.New(op).WithWarpError(err)
	}
	resp, err := s.categoryProvider.GetCategory(ctx, categoryparam.GetCategoryRequest{
		ID: product.CategoryID,
	})
	if err != nil {
		return productparam.UpdateSingleAttributeResponse{}, richerror.New(op).WithWarpError(err)
	}

	validKeys := make(map[string]struct{})
	for _, def := range resp.Category.AttributeSet {
		validKeys[def.Key] = struct{}{}
	}

	if _, ok := validKeys[req.Key]; !ok {
		return productparam.UpdateSingleAttributeResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("product.invalid_specification_key").
			WithMeta(map[string]any{"invalid_key": req.Key})
	}

	if err := s.repo.UpdateSingleSpecification(ctx, req.ProductID, req.Key, req.Value); err != nil {
		return productparam.UpdateSingleAttributeResponse{}, richerror.New(op).WithWarpError(err)
	}

	return productparam.UpdateSingleAttributeResponse{}, nil
}

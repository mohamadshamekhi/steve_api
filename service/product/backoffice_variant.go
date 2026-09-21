package productservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s BackofficeService) AddVariant(ctx context.Context, req productparam.AddVariantRequest) (productparam.AddVariantResponse, error) {
	const op = "productservice.AddVariant"

	v := entity.Variant{
		Price:          req.Price,
		Quantity:       req.Quantity,
		StockMode:      req.StockMode,
		AllowBackorder: req.AllowBackorder,
		BackorderLimit: req.BackorderLimit,
		SupplierID:     req.SupplierID,
		SupplyCost:     req.SupplyCost,
		AvgCost:        req.AvgCost,
		VariantAttribute: []entity.VariantAttribute{
			{
				Key:   req.Key,
				Value: req.Value,
				Label: req.Label,
				Meta:  req.Meta,
			},
		},
		Warranty: req.Warranty,
	}

	// A service holds no inventory, so its variants must never be stock
	// tracked — otherwise selling one would try to draw down a warehouse that
	// does not exist.
	product, pErr := s.repo.FindByID(ctx, req.ProductID)
	if pErr != nil {
		return productparam.AddVariantResponse{}, richerror.New(op).WithWarpError(pErr)
	}
	if product.IsService() {
		// A service is priced once. Letting a second variant in would create a
		// row the service UI never shows and the admin cannot reach, while
		// silently changing which price the invoice picks up.
		if len(product.Variants) > 0 {
			return productparam.AddVariantResponse{}, richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("product.service_single_variant").
				WithMeta(map[string]any{"product_id": req.ProductID})
		}
		v.StockMode = entity.StockModeNone
	}

	if err := s.repo.AddVariant(ctx, req.ProductID, v); err != nil {
		return productparam.AddVariantResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.AddVariantResponse{}, nil
}

func (s BackofficeService) DeleteVariant(ctx context.Context, req productparam.DeleteVariantRequest) (productparam.DeleteVariantResponse, error) {
	const op = "productservice.DeleteVariant"

	// Removing a service's only variant would leave it with no price and make
	// it unsellable in a way the service UI gives no way to repair.
	product, pErr := s.repo.FindByID(ctx, req.ProductID)
	if pErr != nil {
		return productparam.DeleteVariantResponse{}, richerror.New(op).WithWarpError(pErr)
	}
	if product.IsService() {
		return productparam.DeleteVariantResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("product.service_variant_not_deletable").
			WithMeta(map[string]any{"product_id": req.ProductID})
	}

	if err := s.repo.DeleteVariant(ctx, req.ProductID, req.VariantID); err != nil {
		return productparam.DeleteVariantResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.DeleteVariantResponse{}, nil
}

func (s BackofficeService) UpdateVariant(ctx context.Context, req productparam.UpdateVariantRequest) (productparam.UpdateVariantResponse, error) {
	const op = "productservice.UpdateVariant"

	p, err := s.repo.FindByID(ctx, req.ProductID)
	if err != nil {
		return productparam.UpdateVariantResponse{}, richerror.New(op).WithWarpError(err)
	}
	var variant entity.Variant
	for _, v := range p.Variants {
		if v.ID == req.VariantID {
			variant = v
		}
	}
	if variant.ID == "" {
		return productparam.UpdateVariantResponse{}, richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("product.variant_not_found")
	}
	if req.Meta != nil {
		variant.VariantAttribute[0].Meta = *req.Meta
	}
	if req.Key != nil {
		variant.VariantAttribute[0].Key = *req.Key
	}
	if req.Value != nil {
		variant.VariantAttribute[0].Value = *req.Value
	}
	if req.Label != nil {
		variant.VariantAttribute[0].Label = *req.Label
	}
	if req.Price != nil {
		variant.Price = *req.Price
	}
	if req.Quantity != nil {
		variant.Quantity = *req.Quantity
	}
	if req.StockMode != nil {
		variant.StockMode = *req.StockMode
	}
	if req.AllowBackorder != nil {
		variant.AllowBackorder = *req.AllowBackorder
	}
	if req.BackorderLimit != nil {
		variant.BackorderLimit = req.BackorderLimit
	}
	if req.SupplierID != nil {
		variant.SupplierID = *req.SupplierID
	}
	if req.SupplyCost != nil {
		variant.SupplyCost = *req.SupplyCost
	}
	if req.AvgCost != nil {
		variant.AvgCost = *req.AvgCost
	}

	if err := s.repo.UpdateVariant(ctx, req.ProductID, variant); err != nil {
		return productparam.UpdateVariantResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.UpdateVariantResponse{}, nil
}

func (s BackofficeService) UpdateVariantQuantity(ctx context.Context, req productparam.UpdateProductVariantQuantityAfterCancelOrderRequest) (productparam.UpdateProductVariantQuantityAfterCancelOrderResponse, error) {
	const op = "productservice.UpdateVariantQuantity"

	if err := s.repo.UpdateVariantQuantity(ctx, req.ProductID, req.VariantID, req.Quantity); err != nil {
		return productparam.UpdateProductVariantQuantityAfterCancelOrderResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.UpdateProductVariantQuantityAfterCancelOrderResponse{}, nil
}

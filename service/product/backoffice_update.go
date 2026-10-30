package productservice

import (
	"context"
	"errors"
	"fmt"
	"time"

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

		if p.Slug != "" && p.Slug != slug {
			p.PreviousSlugs = appendUnique(p.PreviousSlugs, p.Slug)
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
	if req.RequiresFreightShipping != nil {
		p.RequiresFreightShipping = *req.RequiresFreightShipping
	}
	if req.Kind != nil {
		p.Kind = *req.Kind
	}
	if req.SalesChannels != nil {
		p.SalesChannels = req.SalesChannels
	}
	if req.SeoTitle != nil {
		p.SeoTitle = *req.SeoTitle
	}
	if req.SeoDescription != nil {
		p.SeoDescription = *req.SeoDescription
	}
	if req.FocusKeyword != nil {
		p.FocusKeyword = *req.FocusKeyword
	}
	if req.FAQ != nil {
		p.FAQ = req.FAQ
	}
	if req.TorobURL != nil {
		p.TorobURL = *req.TorobURL
	}

	// Editing a service's fee goes through the product, not the variant: the
	// pricing row is an implementation detail the caller should not have to
	// address. The whole document is replaced below, so mutating it in place
	// keeps the price change atomic with the rest of the update.
	if req.Price != nil && p.IsService() {
		if len(p.Variants) == 0 {
			p.Variants = []entity.Variant{entity.NewServiceVariant(*req.Price)}
			p.Variants[0].BeforeCreate()
			p.Variants[0].StockMode = entity.StockModeNone
		} else {
			p.Variants[0].Price = *req.Price
		}
		// A price change counts as the admin having just reviewed it, so the
		// price-review timer resets instead of flagging this product again
		// right away.
		p.PriceReviewedAt = time.Now()
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return productparam.UpdateProductResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.UpdateProductResponse{Product: productparam.NewProductCreated(p)}, nil
}

// UpdateProductVariantQuantityAfterRegisterOrder draws the sold units out of
// shop stock. It is inventory-aware: services and on-demand variants hold no
// stock at all, and for a backordered line only the part that really came out
// of the shop is deducted — the rest was sourced externally and never sat in
// the warehouse.
//
// The ledger entries for the same sale are written separately by the
// inventory service; this function only maintains the Variant.Quantity cache.
func (s BackofficeService) UpdateProductVariantQuantityAfterRegisterOrder(ctx context.Context, req productparam.UpdateProductVariantQuantityAfterRegisterOrderRequest) (productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse, error) {
	const op = "productservice.UpdateProductVariantQuantityAfterRegisterOrder"

	for _, item := range req.Order.Items {
		if item.GetKind() == entity.ProductKindService {
			continue
		}
		// A hand-typed counter line has no catalog product, and nothing of it
		// ever sat in the warehouse.
		if item.IsCustom || item.ProductID == "" {
			continue
		}

		p, err := s.repo.FindByID(ctx, item.ProductID)
		if err != nil {
			return productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse{}, richerror.New(op).WithWarpError(err)
		}

		var v entity.Variant
		for _, variant := range p.Variants {
			if variant.ID == item.VariantID {
				v = variant
			}
		}
		if v.ID == "" {
			return productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse{}, richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("product.variant_not_found").
				WithMeta(map[string]any{"product_id": item.ProductID, "variant_id": item.VariantID})
		}
		if !v.IsStockTracked() {
			continue
		}

		// Trust the split frozen on the order line when it is there; fall back
		// to recomputing it for orders placed before the feature existed.
		fromStock := uint64(item.StockQuantity)
		if item.Supply == "" {
			fromStock, _ = v.SplitSupply(uint64(item.Quantity))
		}
		if fromStock == 0 {
			continue
		}
		if fromStock > v.Quantity {
			return productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse{}, richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("product.insufficient_stock").
				WithMeta(map[string]any{"product_id": item.ProductID, "variant_id": item.VariantID})
		}

		finalQuantity := uint(v.Quantity - fromStock)
		if err := s.repo.UpdateQuantity(ctx, p.ID, v, finalQuantity); err != nil {
			return productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse{}, richerror.New(op).WithWarpError(err)
		}
	}
	return productparam.UpdateProductVariantQuantityAfterRegisterOrderResponse{}, nil
}

func (s BackofficeService) UpdateProductSalesCountAfterRegisterOrder(ctx context.Context, req productparam.UpdateProductSalesCountAfterRegisterOrderRequest) (productparam.UpdateProductSalesCountAfterRegisterOrderResponse, error) {
	const op = "productservice.UpdateProductSalesCountAfterRegisterOrder"

	for _, item := range req.Order.Items {
		if item.IsCustom || item.ProductID == "" {
			continue
		}
		if err := s.repo.UpdateProductSalesCount(ctx, item.ProductID, item.Quantity); err != nil {
			return productparam.UpdateProductSalesCountAfterRegisterOrderResponse{}, richerror.New(op).WithWarpError(err)
		}
	}
	return productparam.UpdateProductSalesCountAfterRegisterOrderResponse{}, nil
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

func appendUnique(list []string, v string) []string {
	for _, s := range list {
		if s == v {
			return list
		}
	}
	return append(list, v)
}

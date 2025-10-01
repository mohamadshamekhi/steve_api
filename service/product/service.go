package productservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	productparam "github.com/abolfazlnorzad/steve/param/product"
)

type ListResult struct {
	Products   []entity.Product
	TotalCount int64
}

type ProductReader interface {
	FindByID(ctx context.Context, id string) (entity.Product, error)
	FindBySlug(ctx context.Context, slug string) (entity.Product, error)
	All(ctx context.Context, filter commonparam.FilterRequest, page commonparam.PaginationRequest) (ListResult, error)
	GetProductCategoriesForCartItems(ctx context.Context, items []entity.CartItem) (map[string][]string, error)
	GetVariantsByIDs(ctx context.Context, variantIDs []string) (map[string]entity.Variant, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]entity.Product, error)
}

type ProductWriter interface {
	Create(ctx context.Context, product entity.Product) error
	Update(ctx context.Context, product entity.Product) error
	Delete(ctx context.Context, product entity.Product) error
	AddProductGallery(ctx context.Context, productID string, g entity.ProductGallery) error
	DeleteProductGallery(ctx context.Context, productID string, mediaID string) error
	SetPrimary(ctx context.Context, productID string, mediaID string) error
	UpdateSpecifications(ctx context.Context, productID string, specs []entity.Specification) error
	UpdateSingleSpecification(ctx context.Context, productID string, key string, value string) error
	AddVariant(ctx context.Context, ProductID string, v entity.Variant) error
	DeleteVariant(ctx context.Context, ProductID string, VariantID string) error
	UpdateVariant(ctx context.Context, ProductID string, variant entity.Variant) error
	UpdateQuantity(ctx context.Context, ProductID string, variant entity.Variant, finalQuantity uint) error
}

type ProductReadWriter interface {
	ProductReader
	ProductWriter
}

type BackofficeService struct {
	repo             ProductReadWriter
	mediaSvc         MediaService
	categoryProvider CategoryProvider
	brandProvider    BrandProvider
}

func (s BackofficeService) GetVariantsByIDs(ctx context.Context, req productparam.GetVariantsByIDsRequest) (response productparam.GetVariantsByIDsResponse, err error) {
	res, err := s.repo.GetVariantsByIDs(ctx, req.IDs)
	if err != nil {
		return productparam.GetVariantsByIDsResponse{}, err
	}
	return productparam.GetVariantsByIDsResponse{
		Variants: res,
	}, nil
}

func NewBackofficeService(repo ProductReadWriter, mediaSvc MediaService, categoryProvider CategoryProvider, brandProvider BrandProvider) BackofficeService {
	return BackofficeService{repo: repo, mediaSvc: mediaSvc, categoryProvider: categoryProvider, brandProvider: brandProvider}
}

type CategoryProvider interface {
	GetCategory(ctx context.Context, req categoryparam.GetCategoryRequest) (categoryparam.GetCategoryResponse, error)
	GetCategoriesByIDs(ctx context.Context, req categoryparam.GetCategoriesByIDsRequest) (categoryparam.GetCategoriesByIDsResponse, error)
}

type BrandProvider interface {
	GetBrand(ctx context.Context, req brandparam.GetBrandRequest) (brandparam.GetBrandResponse, error)
	GetBrandsByIDs(ctx context.Context, req brandparam.GetBrandsByIDsRequest) (brandparam.GetBrandsByIDsResponse, error)
}

type MediaService interface {
	DeleteMedia(ctx context.Context, req mediaparam.DeleteMediaRequest) (mediaparam.DeleteMediaResponse, error)
	FindMediaByID(ctx context.Context, req mediaparam.GetMediaByIDRequest) (mediaparam.GetMediaByIDResponse, error)
}

type StorefrontService struct {
	repo             ProductReader
	mediaSvc         MediaService
	categoryProvider CategoryProvider
	brandProvider    BrandProvider
}

func NewStorefrontService(
	repo ProductReader,
	mediaSvc MediaService,
	categoryProvider CategoryProvider,
	brandProvider BrandProvider,
) StorefrontService {
	return StorefrontService{
		repo:             repo,
		mediaSvc:         mediaSvc,
		categoryProvider: categoryProvider,
		brandProvider:    brandProvider,
	}
}

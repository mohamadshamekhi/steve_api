package productservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
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
	FindByVariantID(ctx context.Context, variantID string) (entity.Product, error)
	All(ctx context.Context, filter commonparam.FilterRequest, page commonparam.PaginationRequest) (ListResult, error)
	AllForExport(ctx context.Context, filter commonparam.FilterRequest) ([]entity.Product, error)
	GetProductCategoriesForCartItems(ctx context.Context, items []entity.CartItem) (map[string][]string, error)
	GetVariantsByIDs(ctx context.Context, variantIDs []string) (map[string]entity.Variant, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]entity.Product, error)
	FindProductsIDsByCategory(ctx context.Context, categoryID string) (string, error)
	GetFilterableAttributes(ctx context.Context, categoryID string) ([]productparam.AggregatedFilterOption, error)
	GetFilteredProducts(ctx context.Context, filter *productparam.StorefrontFilterProduct, p commonparam.PaginationRequest) (ListResult, error)
	GetFilterMobileListWithMainInfo(ctx context.Context, name *string, brandID *string) ([]entity.Product, error)
	GetMinAdnMaxPriceOfProductWithSpecialCategory(ctx context.Context, categoryID string) (uint64, uint64, error)
	UpdateProductAvrRating(ctx context.Context, ProductID string, AverageRating float64) error
	GetProductIDsByCategories(ctx context.Context, catIDs []string) (map[string][]string, error)
	GetPricesOfProducts(ctx context.Context, prodIDSlice []string, variantIDSlice []string) (map[string]map[string]float64, error)
	GetSearchProduct(ctx context.Context, search commonparam.SearchRequest) (ListResult, error)
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
	UpdateProductSalesCount(ctx context.Context, ProductID string, quantity uint) error
	UpdateVariantQuantity(ctx context.Context, productID string, variantID string, quantity uint) error
	MarkPriceReviewed(ctx context.Context, productID string, at time.Time) error
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
	GetCategoryBySlug(ctx context.Context, req categoryparam.GetCategoryBySlugRequest) (categoryparam.GetCategoryBySlugResponse, error)
}

type BrandProvider interface {
	GetBrand(ctx context.Context, req brandparam.GetBrandRequest) (brandparam.GetBrandResponse, error)
	GetBrandsByIDs(ctx context.Context, req brandparam.GetBrandsByIDsRequest) (brandparam.GetBrandsByIDsResponse, error)
}

type MediaService interface {
	DeleteMedia(ctx context.Context, req mediaparam.DeleteMediaRequest) (mediaparam.DeleteMediaResponse, error)
	FindMediaByID(ctx context.Context, req mediaparam.GetMediaByIDRequest) (mediaparam.GetMediaByIDResponse, error)
}

type DiscountProvider interface {
	GetProductIDsThatHaveActiveAutoDiscount(ctx context.Context, _ discountparam.GetProductIDsThatHaveActiveAutoDiscountRequest) (discountparam.GetProductIDsThatHaveActiveAutoDiscountResponse, error)
	GetSpecialDiscount(ctx context.Context, req discountparam.GetSpecialDiscountRequest) (discountparam.GetSpecialDiscountResponse, error)
}

type StorefrontService struct {
	repo             ProductReadWriter
	mediaSvc         MediaService
	categoryProvider CategoryProvider
	brandProvider    BrandProvider
	discountProvider DiscountProvider
}

func NewStorefrontService(
	repo ProductReadWriter,
	mediaSvc MediaService,
	categoryProvider CategoryProvider,
	brandProvider BrandProvider,
	discountProvider DiscountProvider,
) StorefrontService {
	return StorefrontService{
		repo:             repo,
		mediaSvc:         mediaSvc,
		categoryProvider: categoryProvider,
		brandProvider:    brandProvider,
		discountProvider: discountProvider,
	}
}

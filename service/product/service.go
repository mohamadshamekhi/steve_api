package productservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
)

type ListResult struct {
	Products   []entity.Product
	TotalCount int64
}

type ProductReader interface {
	FindByID(ctx context.Context, id string) (entity.Product, error)
	FindBySlug(ctx context.Context, slug string) (entity.Product, error)
	All(ctx context.Context, filter commonparam.FilterRequest, page commonparam.PaginationRequest) (ListResult, error)
}

type ProductWriter interface {
	Create(ctx context.Context, product entity.Product) error
	Update(ctx context.Context, product entity.Product) error
	Delete(ctx context.Context, product entity.Product) error
	AddProductGallery(ctx context.Context, productID string, g entity.ProductGallery) error
	DeleteProductGallery(ctx context.Context, productID string, mediaID string) error
	SetPrimary(ctx context.Context, productID string, mediaID string) error
}

type ProductReadWriter interface {
	ProductReader
	ProductWriter
}

type BackofficeService struct {
	repo     ProductReadWriter
	mediaSvc MediaService
}

func NewBackofficeService(repo ProductReadWriter, mediaSvc MediaService) BackofficeService {
	return BackofficeService{repo: repo, mediaSvc: mediaSvc}
}

type MediaService interface {
	DeleteMedia(ctx context.Context, req mediaparam.DeleteMediaRequest) (mediaparam.DeleteMediaResponse, error)
	FindMediaByID(ctx context.Context, req mediaparam.GetMediaByIDRequest) (mediaparam.GetMediaByIDResponse, error)
}

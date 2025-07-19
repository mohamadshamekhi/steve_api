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
}

type ArticleReadWriter interface {
	ProductReader
	ProductWriter
}

type BackofficeService struct {
	repo ArticleReadWriter
}

type MediaService interface {
	GetURLByMediaID(ctx context.Context, req mediaparam.GetUrlByIDRequest) (mediaparam.GetUrlByIDResponse, error)
	GetURLsByMediaIDs(ctx context.Context, req mediaparam.GetUrlsByIDsRequest) (mediaparam.GetUrlsByIDsResponse, error)
	DeleteMedia(ctx context.Context, req mediaparam.DeleteMediaRequest) (mediaparam.DeleteMediaResponse, error)
	FindMediaByID(ctx context.Context, req mediaparam.GetMediaByIDRequest) (mediaparam.GetMediaByIDResponse, error)
}

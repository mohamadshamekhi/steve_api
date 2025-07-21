package categoryservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
)

type ListResult struct {
	Categories []entity.Category
	TotalCount int64
}

type CategoryReader interface {
	FindByID(ctx context.Context, id string) (entity.Category, error)
	FindBySlug(ctx context.Context, slug string) (entity.Category, error)
	All(ctx context.Context, filter commonparam.FilterRequest, page commonparam.PaginationRequest) (ListResult, error)
}

type CategoryWriter interface {
	Create(ctx context.Context, category entity.Category) error
	Update(ctx context.Context, category entity.Category) error
	Delete(ctx context.Context, category entity.Category) error
	GetCategoriesByIDs(ctx context.Context, ids []string) (map[string]entity.Category, error)
}

type CategoryReadWriter interface {
	CategoryReader
	CategoryWriter
}
type MediaService interface {
	GetURLByMediaID(ctx context.Context, req mediaparam.GetUrlByIDRequest) (mediaparam.GetUrlByIDResponse, error)
	GetURLsByMediaIDs(ctx context.Context, req mediaparam.GetUrlsByIDsRequest) (mediaparam.GetUrlsByIDsResponse, error)
	DeleteMedia(ctx context.Context, req mediaparam.DeleteMediaRequest) (mediaparam.DeleteMediaResponse, error)
	FindMediaByID(ctx context.Context, req mediaparam.GetMediaByIDRequest) (mediaparam.GetMediaByIDResponse, error)
}

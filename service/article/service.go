package articleservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
)

type ListResult struct {
	Articles   []entity.Article
	TotalCount int64
}

type ArticleReader interface {
	FindByID(ctx context.Context, id string) (entity.Article, error)
	FindBySlug(ctx context.Context, slug string) (entity.Article, error)
	All(ctx context.Context, filter commonparam.FilterRequest, page commonparam.PaginationRequest) (ListResult, error)
}

type ArticleWriter interface {
	Create(ctx context.Context, category entity.Article) error
	Update(ctx context.Context, category entity.Article) error
	Delete(ctx context.Context, category entity.Article) error
}

type ArticleReadWriter interface {
	ArticleReader
	ArticleWriter
}
type MediaService interface {
	GetURLByMediaID(ctx context.Context, req mediaparam.GetUrlByIDRequest) (mediaparam.GetUrlByIDResponse, error)
	GetURLsByMediaIDs(ctx context.Context, req mediaparam.GetUrlsByIDsRequest) (mediaparam.GetUrlsByIDsResponse, error)
	DeleteMedia(ctx context.Context, req mediaparam.DeleteMediaRequest) (mediaparam.DeleteMediaResponse, error)
	FindMediaByID(ctx context.Context, req mediaparam.GetMediaByIDRequest) (mediaparam.GetMediaByIDResponse, error)
}

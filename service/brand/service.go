package brandservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
)

type ListResult struct {
	Brands     []entity.Brand
	TotalCount int64
}

type BrandReader interface {
	FindByID(ctx context.Context, id string) (entity.Brand, error)
	All(ctx context.Context, filter commonparam.FilterRequest, page commonparam.PaginationRequest) (ListResult, error)
}

type BrandWriter interface {
	Create(ctx context.Context, brand entity.Brand) error
	Update(ctx context.Context, brand entity.Brand) error
	Delete(ctx context.Context, brand entity.Brand) error
}

type BrandReadWriter interface {
	BrandReader
	BrandWriter
}

type MediaService interface {
	GetURLByMediaID(ctx context.Context, req mediaparam.GetUrlByIDRequest) (mediaparam.GetUrlByIDResponse, error)
	GetURLsByMediaIDs(ctx context.Context, req mediaparam.GetUrlsByIDsRequest) (mediaparam.GetUrlsByIDsResponse, error)
	DeleteMedia(ctx context.Context, req mediaparam.DeleteMediaRequest) (mediaparam.DeleteMediaResponse, error)
	FindMediaByID(ctx context.Context, req mediaparam.GetMediaByIDRequest) (mediaparam.GetMediaByIDResponse, error)
}

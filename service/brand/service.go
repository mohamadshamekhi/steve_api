package brandservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
)

type BrandReader interface {
	FindByID(ctx context.Context, id string) (entity.Brand, error)
	//ListAll(ctx context.Context) ([]entity.Brand, error)
}

type BrandWriter interface {
	Create(ctx context.Context, brand entity.Brand) error
	Update(ctx context.Context, brand entity.Brand) error
}

type BrandReadWriter interface {
	BrandReader
	BrandWriter
}

type MediaService interface {
	GetURLByMediaID(ctx context.Context, req mediaparam.GetUrlByIDRequest) (mediaparam.GetUrlByIDResponse, error)
	DeleteMedia(ctx context.Context, req mediaparam.DeleteMediaRequest) (mediaparam.DeleteMediaResponse, error)
}

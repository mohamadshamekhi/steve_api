package mediaservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"mime/multipart"
)

type ObjectStorageClient interface {
	Upload(ctx context.Context, header multipart.FileHeader, filePath string) error
	GetURL(ctx context.Context, filePath string) (string, error)
	Delete(ctx context.Context, filePath string) error
}

type Repository interface {
	AddMedia(ctx context.Context, media entity.Media) (entity.Media, error)
	GetMediaByID(ctx context.Context, id uint) (entity.Media, error)
	DeleteMedia(ctx context.Context, id uint) error
}
type Service struct {
	client     ObjectStorageClient
	repository Repository
}

func New(client ObjectStorageClient, r Repository) Service {
	return Service{client: client, repository: r}
}

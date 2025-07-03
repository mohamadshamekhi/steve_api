package mediaservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/cache"
	"mime/multipart"
	"time"
)

type Config struct {
	PublicBaseURL                 string        `koanf:"public_base_url"`
	AllowedExtensions             []string      `koanf:"allowed_extensions"`
	AllowedMimeTypes              []string      `koanf:"allowed_mime_types"`
	MaxUploadSizeBytes            int64         `koanf:"max_upload_size"`
	PresignedPrivateURLExpiration time.Duration `koanf:"presigned_private_url_expiration"`
	PresignedPublicURLExpiration  time.Duration `koanf:"presigned_public_url_expiration"`
}

type UploadResult struct {
	Path       string
	BucketName string
	MimeType   string
}
type StorageProvider interface {
	Upload(ctx context.Context, domain string, isPublic bool, file multipart.File, header *multipart.FileHeader) (UploadResult, error)
	GetURL(ctx context.Context, media *entity.Media) (string, error)
	Delete(ctx context.Context, media *entity.Media) error
}

type Repository interface {
	AddMedia(ctx context.Context, media entity.Media) (entity.Media, error)
	GetMediaByID(ctx context.Context, id string) (entity.Media, error)
	DeleteMedia(ctx context.Context, id string) error
}
type Service struct {
	provider   StorageProvider
	repository Repository
	config     Config
	cache      cache.Provider
}

func New(provider StorageProvider, r Repository, cfg Config, cache cache.Provider) Service {
	return Service{provider: provider, repository: r, config: cfg, cache: cache}
}

package mediaservice

import (
	"context"
	"io"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/cache"
)

type Config struct {
	PublicBaseURL                 string                 `koanf:"public_base_url"`
	AllowedExtensions             []string               `koanf:"allowed_extensions"`
	AllowedMimeTypes              []string               `koanf:"allowed_mime_types"`
	MaxUploadSizeBytes            int64                  `koanf:"max_upload_size"`
	PresignedPrivateURLExpiration time.Duration          `koanf:"presigned_private_url_expiration"`
	PresignedPublicURLExpiration  time.Duration          `koanf:"presigned_public_url_expiration"`
	Policies                      map[string]MediaPolicy `koanf:"policies"`
}

type MediaPolicy struct {
	MaxUploadSizeBytes int64    `koanf:"max_upload_size"`
	AllowedMimeTypes   []string `koanf:"allowed_mime_types"`
}

type UploadResult struct {
	Path       string
	BucketName string
	MimeType   string
}
type StorageProvider interface {
	Upload(ctx context.Context, domain string, role string, isPublic bool, file io.Reader, info mediaparam.FileMetadata) (UploadResult, error)
	GetURL(ctx context.Context, media *entity.Media) (string, error)
	Delete(ctx context.Context, media *entity.Media) error
	TransformToWebP(input io.Reader) (io.ReadSeeker, int64, error)
}

type Repository interface {
	AddMedia(ctx context.Context, media entity.Media) (entity.Media, error)
	GetMediaByID(ctx context.Context, id string) (entity.Media, error)
	DeleteMedia(ctx context.Context, id string) error
	FindManyByIDs(ctx context.Context, ids []string) ([]entity.Media, error)
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

func (c *Config) GetPolicy(role string) MediaPolicy {
	if p, ok := c.Policies[role]; ok {
		return p
	}

	return MediaPolicy{
		MaxUploadSizeBytes: c.MaxUploadSizeBytes,
		AllowedMimeTypes:   c.AllowedMimeTypes,
	}
}

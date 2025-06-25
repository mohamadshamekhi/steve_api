package entity

import (
	"github.com/oklog/ulid/v2"
	"time"
)

var ValidExt = []string{"mp4", "png", "jpg", "jpeg", "webp", "avi", "webm"}

const (
	MaxFileUploadSize          = 300 * 1024 * 1024
	FileLinkExpirationDuration = 60 * int64(time.Second) * 60 * 7
)

type Media struct {
	ID              string                 `bson:"_id"`
	StorageProvider string                 `bson:"storage_provider"` // "s3", "gcs", "local"
	BucketName      string                 `bson:"bucket_name"`
	Path            string                 `bson:"path"`
	IsPublic        bool                   `bson:"is_public"`
	MimeType        string                 `bson:"mime_type"`
	Size            int64                  `bson:"size_in_bytes"`
	UploaderID      string                 `bson:"uploader_id"`
	Meta            map[string]interface{} `bson:"meta,omitempty"`
	CreatedAt       time.Time              `bson:"created_at"`
	UpdatedAt       time.Time              `bson:"updated_at"`
	DeletedAt       time.Time              `bson:"deleted_at,omitempty"`
}

func (m *Media) BeforeCreate() {
	m.ID = ulid.Make().String()
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
}

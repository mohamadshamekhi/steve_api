package entity

import (
	"github.com/oklog/ulid/v2"
	"time"
)

type Media struct {
	ID               string                 `bson:"_id"`
	OriginalFileName string                 `bson:"original_file_name"`
	StorageProvider  string                 `bson:"storage_provider"` // "s3", "gcs", "local"
	Path             string                 `bson:"path"`
	IsPublic         bool                   `bson:"is_public"`
	MimeType         string                 `bson:"mime_type"`
	Size             int64                  `bson:"size_in_bytes"`
	UploaderID       string                 `bson:"uploader_id"`
	UploaderType     string                 `bson:"uploader_type"`
	Meta             map[string]interface{} `bson:"meta,omitempty"`
	CreatedAt        time.Time              `bson:"created_at"`
	UpdatedAt        time.Time              `bson:"updated_at"`
	DeletedAt        time.Time              `bson:"deleted_at,omitempty"`
}

func (m *Media) BeforeCreate() {
	m.ID = ulid.Make().String()
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
}

package entity

import "time"

var ValidExt = []string{"mp4", "png", "jpg", "jpeg", "webp", "avi", "webm"}

const (
	MaxFileUploadSize          = 300 * 1024 * 1024
	FileLinkExpirationDuration = 60 * int64(time.Second) * 60 * 7
)

type Media struct {
	ID          string    `bson:"_id"`
	Kind        string    `bson:"kind"`
	FileName    string    `bson:"filename"`
	SubDir      string    `bson:"subdir"`
	StorageSize uint      `bson:"storage_size"`
	CreatedAt   time.Time `bson:"created_at"`
}

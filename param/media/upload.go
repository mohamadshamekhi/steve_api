package mediaparam

import "mime/multipart"

type UploadMediaRequest struct {
	Domain        string                `json:"domain"`
	UploaderID    string                `json:"uploader_id"`
	UploaderType  string                `json:"uploader_type"`
	IsPublic      bool                  `json:"is_public"`
	File          multipart.File        `json:"file"`
	Header        *multipart.FileHeader `json:"header"`
	ConvertToWebp bool                  `json:"convert_to_webp"`
}

type UploadMediaResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
type FileMetadata struct {
	Filename    string `json:"filename"`
	Size        uint64 `json:"size"`
	ContentType string `json:"content_type"`
}

type StorefrontUploadMediaRequest struct {
	Domain        string `json:"domain"`
	UploaderID    string `json:"uploader_id"`
	UploaderType  string
	File          multipart.File        `json:"file"`
	Header        *multipart.FileHeader `json:"header"`
	ConvertToWebp bool                  `json:"convert_to_webp"`
}

type StorefrontUploadMediaResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

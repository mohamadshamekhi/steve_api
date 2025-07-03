package mediaparam

import "mime/multipart"

type UploadMediaRequest struct {
	Domain       string                `json:"domain"`
	UploaderID   string                `json:"uploader_id"`
	UploaderType string                `json:"uploader_type"`
	IsPublic     bool                  `json:"is_public"`
	File         multipart.File        `json:"file"`
	Header       *multipart.FileHeader `json:"header"`
}

type UploadMediaResponse struct {
	ID string `json:"id"`
}

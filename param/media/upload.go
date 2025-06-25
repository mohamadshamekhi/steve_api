package mediaparam

import "mime/multipart"

type UploadMediaRequest struct {
	*multipart.FileHeader `json:"file"`
	MediaableID           uint   `json:"mediaable_id"`
	MediaableType         string `json:"mediaable_type"`
}

type UploadMediaResponse struct {
	ID       uint   `json:"id"`
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
}

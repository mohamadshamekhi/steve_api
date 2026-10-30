package articleparam

type ExternalizeInlineImagesRequest struct{}

type ExternalizedArticle struct {
	ID              string `json:"id"`
	Slug            string `json:"slug"`
	Images          int    `json:"images"`
	BodyBytesBefore int    `json:"body_bytes_before"`
	BodyBytesAfter  int    `json:"body_bytes_after"`
}

type ExternalizeInlineImagesResponse struct {
	Fixed  []ExternalizedArticle `json:"fixed"`
	Failed []ExternalizedArticle `json:"failed"`
}

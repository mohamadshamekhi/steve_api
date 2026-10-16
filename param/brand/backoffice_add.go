package brandparam

type CreateBrandRequest struct {
	Name string `json:"name"`
	// Slug is optional. When empty it's auto-generated from Name (like
	// category); pass it explicitly to pin an English slug (e.g. "samsung")
	// for a brand whose Name is stored in Persian.
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	MediaID     string `json:"media_id,omitempty"`
	IsFavorite  bool   `json:"is_favorite"`
}

type CreateBrandResponse struct {
}

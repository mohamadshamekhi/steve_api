package brandparam

type CreateBrandRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MediaID     string `json:"media_id,omitempty"`
	IsFavorite  bool   `json:"is_favorite"`
}

type CreateBrandResponse struct {
}

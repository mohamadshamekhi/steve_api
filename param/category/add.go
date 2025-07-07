package categoryparam

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MediaID     string `json:"media_id,omitempty"`
	ParentID    string `json:"parent_id,omitempty"`
}

type CreateCategoryResponse struct {
}

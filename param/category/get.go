package categoryparam

type GetCategoryRequest struct {
	ID string `json:"id"`
}

type GetCategoryResponse struct {
	Category CategoryInfo `json:"category"`
}

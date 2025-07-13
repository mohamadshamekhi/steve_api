package categoryparam

type UpdateCategoryRequest struct {
	CategoryID  string  `json:"category_id"`
	Name        *string `json:"name"`
	MediaID     *string `json:"media_id"`
	Description *string `json:"description"`
}
type UpdateCategoryResponse struct{}

package categoryparam

type UpdateCategoryRequest struct {
	CategoryID     string  `json:"category_id"`
	Name           *string `json:"name"`
	MediaID        *string `json:"media_id"`
	Description    *string `json:"description"`
	ParentID       *string `json:"parent_id"`
	AttributeSetID *string `json:"attribute_set_id"`
}
type UpdateCategoryResponse struct{}

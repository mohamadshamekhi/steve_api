package categoryparam

type GetFilterForCategoryRequest struct {
	CategoryID string `json:"category_id"`
}

type GetFilterForCategoryResponse struct {
	FilterOption []FilterOption `json:"filter_option"`
}
type FilterOption struct {
	Key     string   `json:"key"`
	Name    string   `json:"name"`
	Unit    string   `json:"unit,omitempty"`
	Options []string `json:"options"`
}

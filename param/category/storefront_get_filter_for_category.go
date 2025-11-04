package categoryparam

type GetFilterForCategoryRequest struct {
	CategoryID string `json:"category_id"`
}

type GetFilterForCategoryResponse struct {
	FilterOption []FilterOption `json:"filter_option"`
	MaxPrice     uint64         `json:"max_price"`
	MinPrice     uint64         `json:"min_price"`
}
type FilterOption struct {
	Key     string   `json:"key"`
	Name    string   `json:"name"`
	Unit    string   `json:"unit,omitempty"`
	Options []string `json:"options"`
}

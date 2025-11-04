package productparam

type GetMinAdnMaxPriceOfProductWithSpecialCategoryRequest struct {
	CategoryID string `json:"category_id"`
}
type GetMinAdnMaxPriceOfProductWithSpecialCategoryResponse struct {
	MaxPrice uint64 `json:"max_price"`
	MinPrice uint64 `json:"min_price"`
}

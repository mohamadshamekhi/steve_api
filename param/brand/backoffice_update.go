package brandparam

type UpdateBrandRequest struct {
	BrandID     string  `json:"brand_id"`
	Name        *string `json:"name"`
	MediaID     *string `json:"media_id"`
	Description *string `json:"description"`
}
type UpdateBrandResponse struct{}

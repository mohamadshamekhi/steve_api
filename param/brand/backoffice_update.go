package brandparam

type UpdateBrandRequest struct {
	BrandID string  `json:"brand_id"`
	Name    *string `json:"name"`
	// Slug is optional; when set (non-empty) it overrides the auto-generated
	// slug. Leaving it nil/empty while changing Name regenerates the slug
	// from the new Name, same as category.
	Slug        *string `json:"slug"`
	MediaID     *string `json:"media_id"`
	Description *string `json:"description"`
	IsFavorite  *bool   `json:"is_favorite"`
}
type UpdateBrandResponse struct{}

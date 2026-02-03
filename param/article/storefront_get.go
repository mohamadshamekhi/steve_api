package articleparam

type GetStorefrontArticleRequest struct {
	ID string `json:"id"`
}
type GetStorefrontArticleResponse struct {
	Article StorefrontArticleDetail `json:"article"`
}

type GetBySlugStorefrontArticleRequest struct {
	Slug string `json:"slug"`
}
type GetBySlugStorefrontArticleResponse struct {
	Article StorefrontArticleDetail `json:"article"`
}

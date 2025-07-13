package articleparam

type GetStorefrontArticleRequest struct {
	ID string `json:"id"`
}
type GetStorefrontArticleResponse struct {
	Article StorefrontArticleDetail `json:"article"`
}

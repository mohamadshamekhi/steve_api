package articleparam

type GetStorefrontArticleRequest struct {
	ID string `json:"id"`
	// VisitorID identifies the browser making the request (set from an
	// anonymous cookie) so a repeat visit doesn't get double-counted. Never
	// part of the request body/query - filled in by the handler.
	VisitorID string `json:"-"`
}
type GetStorefrontArticleResponse struct {
	Article StorefrontArticleDetail `json:"article"`
}

type GetBySlugStorefrontArticleRequest struct {
	Slug      string `json:"slug"`
	VisitorID string `json:"-"`
}
type GetBySlugStorefrontArticleResponse struct {
	Article StorefrontArticleDetail `json:"article"`
}

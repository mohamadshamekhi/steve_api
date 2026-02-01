package articleparam

type UpdateArticleRequest struct {
	ID              string    `json:"id"`
	Title           *string   `json:"title"`
	Slug            *string   `json:"slug"`
	Body            *string   `json:"body"`
	MediaID         *string   `json:"media_id"`
	IsPublished     *bool     `json:"is_published"`
	Tags            *[]string `json:"tags"`
	MetaDescription *string   `json:"meta_description"`
}

type UpdateArticleResponse struct {
	Article ArticleInfo `json:"article"`
}

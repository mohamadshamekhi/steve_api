package articleparam

type UpdateArticleRequest struct {
	ID          string  `json:"id"`
	Title       *string `json:"title"`
	Body        *string `json:"body"`
	MediaID     *string `json:"media_id"`
	IsPublished bool    `json:"is_published"`
}

type UpdateArticleResponse struct {
	Article ArticleInfo `json:"article"`
}

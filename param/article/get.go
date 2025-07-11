package articleparam

type GetArticleRequest struct {
	ID string `json:"id"`
}
type GetArticleResponse struct {
	Article ArticleInfo `json:"article"`
}

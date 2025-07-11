package articleparam

import commonparam "github.com/abolfazlnorzad/steve/param/common"

type AllArticleRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}

type AllArticleResponse struct {
	Articles   []ArticleInfo                  `json:"articles"`
	Pagination commonparam.PaginationResponse `json:"pagination"`
}

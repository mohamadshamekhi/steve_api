package commentparam

import (
	"time"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllCommentsRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllCommentsResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Comments   []CommentsInfo                 `json:"comments"`
}

type CommentsInfo struct {
	ID         string    `json:"id"`
	Body       string    `json:"body"`
	IsApproved bool      `json:"is_approved"`
	CreatedAt  time.Time `json:"created_at"`

	Commenter struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"` // "user" or "backoffice"
	} `json:"commenter"`

	Commentable struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"` // "product" or "article"
	} `json:"commentable"`

	ParentComment *struct {
		ID            string `json:"id"`
		Body          string `json:"body"`
		CommenterName string `json:"commenter_name"`
	} `json:"parent_comment,omitempty"`
}

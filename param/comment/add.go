package commentparam

type AddCommentRequest struct {
	Body            string  `json:"body"`
	CommentableID   string  `json:"commentable_id"`
	CommentableType string  `json:"commentable_type"`
	CommenterName   string  `json:"commenter_name"`
	CommenterType   string  `json:"commenter_type"`
	CommenterID     string  `json:"commenter_id"`
	ParentID        *string `json:"parent_id,omitempty"`
}

type AddCommentResponse struct{}

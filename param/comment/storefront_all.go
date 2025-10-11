package commentparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type StorefrontGetCommentsRequest struct {
	Pagination      commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	CommentableType string
	CommentableID   string
}

type StorefrontGetCommentsResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Comments   []CommentWithChildren          `json:"comments"`
}

type CommentWithChildren struct {
	ID              string                   `bson:"_id" json:"id"`
	Body            string                   `bson:"body" json:"body"`
	CommentableID   string                   `bson:"commentable_id" json:"commentable_id"`
	CommentableType string                   `bson:"commentable_type" json:"commentable_type"`
	CommenterID     string                   `bson:"commenter_id" json:"commenter_id"`
	CommenterType   string                   `bson:"commenter_type" json:"commenter_type"`
	CommenterName   string                   `bson:"commenter_name" json:"commenter_name"`
	ParentID        string                   `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Ancestors       []entity.AncestorComment `bson:"ancestors,omitempty" json:"ancestors,omitempty"`
	Level           int                      `bson:"level" json:"level"`
	IsApproved      bool                     `bson:"is_approved" json:"is_approved"`
	CreatedAt       time.Time                `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time                `bson:"updated_at" json:"updated_at"`
	Children        []*CommentWithChildren   `bson:"children" json:"children"`
}

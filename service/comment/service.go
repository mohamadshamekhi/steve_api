package commentservice

import (
	"context"
	"errors"
	"github.com/abolfazlnorzad/steve/entity"
	commentparam "github.com/abolfazlnorzad/steve/param/comment"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (entity.Comment, error)
	Create(ctx context.Context, comment entity.Comment) error
}
type Service struct {
	repo Repository
}

func (s Service) AddComment(ctx context.Context, req commentparam.AddCommentRequest) (commentparam.AddCommentResponse, error) {
	const op = "commentservice.AddComment"
	var pid string
	if req.ParentID != nil {
		pid = *req.ParentID
	}
	cm := entity.Comment{
		Body:            req.Body,
		CommentableID:   req.CommentableID,
		CommentableType: req.CommentableType,
		CommenterID:     req.CommentableID,
		CommenterType:   req.CommentableType,
		CommenterName:   req.CommenterName,
		ParentID:        pid,
	}
	if req.ParentID != nil {
		p, fErr := s.repo.FindByID(ctx, *req.ParentID)
		if fErr != nil {
			var rErr *richerror.RichError
			if errors.As(fErr, &rErr) && rErr.Kind() == richerror.KindNotFound {
				return commentparam.AddCommentResponse{}, richerror.New(op).
					WithKind(richerror.KindInvalid).
					WithMessageKey("global.not-found").
					WithMeta(map[string]any{"parent_id": req.ParentID})
			}
			return commentparam.AddCommentResponse{}, richerror.New(op).WithWarpError(fErr)
		}

		cm.Ancestors = p.Ancestors
		cm.Ancestors = append(cm.Ancestors, entity.AncestorComment{
			ID:            p.ID,
			CommenterID:   p.CommenterID,
			CommenterType: p.CommenterType,
			CommenterName: p.CommenterName,
			Body:          p.Body,
			Level:         p.Level,
		})
		cm.Level = p.Level + 1
	} else {
		cm.Level = 0
		cm.Ancestors = []entity.AncestorComment{}
	}
	cm.BeforeCreate()
	if err := s.repo.Create(ctx, cm); err != nil {
		return commentparam.AddCommentResponse{}, richerror.New(op).WithWarpError(err)
	}
	return commentparam.AddCommentResponse{}, nil
}

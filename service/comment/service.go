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
	Update(ctx context.Context, comment entity.Comment) error
}
type Service struct {
	repo              Repository
	commenterFetchers map[string]CommenterFetcher
}

type Commenter interface {
	GetID() string
	GetName() string
}

type CommenterFetcher interface {
	FindCommenterByID(ctx context.Context, id string) (Commenter, error)
}

func New(repo Repository, fetchers map[string]CommenterFetcher) Service {
	return Service{repo: repo, commenterFetchers: fetchers}
}

func (s Service) AddComment(ctx context.Context, req commentparam.AddCommentRequest) (commentparam.AddCommentResponse, error) {
	const op = "commentservice.AddComment"
	var pid string
	if req.ParentID != nil {
		pid = *req.ParentID
	}
	fetcher, ok := s.commenterFetchers[req.CommenterType]
	if !ok {
		return commentparam.AddCommentResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("comment.invalid_commenter_type").
			WithMeta(map[string]any{"type": req.CommenterType})
	}
	commenter, err := fetcher.FindCommenterByID(ctx, req.CommenterID)
	if err != nil {
		return commentparam.AddCommentResponse{}, richerror.New(op).WithWarpError(err)
	}
	cm := entity.Comment{
		Body:            req.Body,
		CommentableID:   req.CommentableID,
		CommentableType: req.CommentableType,
		CommenterID:     commenter.GetID(),
		CommenterType:   req.CommenterType,
		CommenterName:   commenter.GetName(),
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

func (s Service) UpdateComment(ctx context.Context, req commentparam.UpdateCommentRequest) (commentparam.UpdateCommentResponse, error) {
	const op = "commentservice.UpdateComment"
	cm, err := s.repo.FindByID(ctx, req.CommentID)
	if err != nil {
		return commentparam.UpdateCommentResponse{}, richerror.New(op).WithWarpError(err)
	}
	cm.IsApproved = !cm.IsApproved
	if err := s.repo.Update(ctx, cm); err != nil {
		return commentparam.UpdateCommentResponse{}, richerror.New(op).WithWarpError(err)
	}
	return commentparam.UpdateCommentResponse{}, nil
}

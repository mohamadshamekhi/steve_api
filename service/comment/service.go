package commentservice

import (
	"context"
	"errors"
	"sync"

	"github.com/abolfazlnorzad/steve/entity"
	commentparam "github.com/abolfazlnorzad/steve/param/comment"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type ListResult struct {
	Comments   []entity.Comment
	TotalCount uint
}
type StorefrontListResult struct {
	Comments   []commentparam.CommentWithChildren
	TotalCount uint
}

type Repository interface {
	FindByID(ctx context.Context, id string) (entity.Comment, error)
	Create(ctx context.Context, comment entity.Comment) error
	Update(ctx context.Context, comment entity.Comment) error

	List(ctx context.Context, filter commonparam.FilterRequest, pagination commonparam.PaginationRequest) (ListResult, error)
	FindManyByIDs(ctx context.Context, ids []string) (map[string]entity.Comment, error)
	GetApprovedCommentsWithChildren(ctx context.Context, commentableType string, commentableID string, pagination commonparam.PaginationRequest) (StorefrontListResult, error)
}
type Service struct {
	repo               Repository
	commentableFetcher CommentableFetcher
	commenterFetcher   CommenterFetcher
}

type Commentable interface {
	GetID() string
	GetTitle() string
}

type CommentableFetcher interface {
	FindCommentablesByIDs(ctx context.Context, idsByType map[string][]string) (map[string]Commentable, error)
}

type Commenter interface {
	GetID() string
	GetName() string
}

type CommenterFetcher interface {
	FindCommentersByIDs(ctx context.Context, idsByType map[string][]string) (map[string]Commenter, error)
}

func New(repo Repository, fetchers CommenterFetcher, commentableFetcher CommentableFetcher) Service {
	return Service{repo: repo, commenterFetcher: fetchers, commentableFetcher: commentableFetcher}
}

func (s Service) AddComment(ctx context.Context, req commentparam.AddCommentRequest) (commentparam.AddCommentResponse, error) {
	const op = "commentservice.AddComment"
	commenterMap, err := s.commenterFetcher.FindCommentersByIDs(ctx, map[string][]string{req.CommenterType: {req.CommenterID}})
	if err != nil {
		return commentparam.AddCommentResponse{}, richerror.New(op).WithWarpError(err)
	}
	commenter, ok := commenterMap[req.CommenterID]
	if !ok {
		return commentparam.AddCommentResponse{}, richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("user.not_found")
	}

	newComment := &entity.Comment{
		Body:            req.Body,
		CommentableID:   req.CommentableID,
		CommentableType: req.CommentableType,
		CommenterID:     commenter.GetID(),
		CommenterType:   req.CommenterType,
		CommenterName:   commenter.GetName(),
	}

	if req.ParentID != nil && *req.ParentID != "" {
		parent, err := s.repo.FindByID(ctx, *req.ParentID)
		if err != nil {
			var rErr *richerror.RichError
			if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
				return commentparam.AddCommentResponse{}, richerror.New(op).
					WithKind(richerror.KindInvalid).
					WithMessageKey("comment.parent_not_found").
					WithMeta(map[string]any{"parent_id": *req.ParentID})
			}
			return commentparam.AddCommentResponse{}, richerror.New(op).WithWarpError(err)
		}

		newComment.ParentID = parent.ID
		newComment.Ancestors = parent.Ancestors
		newComment.Ancestors = append(newComment.Ancestors, entity.AncestorComment{
			ID:            parent.ID,
			CommenterID:   parent.CommenterID,
			CommenterType: parent.CommenterType,
			CommenterName: parent.CommenterName,
			Body:          parent.Body,
			Level:         parent.Level,
		})
		newComment.Level = parent.Level + 1
	} else {
		newComment.Level = 0
		newComment.Ancestors = []entity.AncestorComment{}
	}

	newComment.BeforeCreate()
	if err := s.repo.Create(ctx, *newComment); err != nil {
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

func (s Service) AllComments(ctx context.Context, req commentparam.GetAllCommentsRequest) (commentparam.GetAllCommentsResponse, error) {
	const op = "commentservice.AllComments"
	listResult, err := s.repo.List(ctx, req.Filter, req.Pagination)
	if err != nil {
		return commentparam.GetAllCommentsResponse{}, richerror.New(op).WithWarpError(err)
	}

	if len(listResult.Comments) == 0 {
		return commentparam.GetAllCommentsResponse{
			Pagination: commonparam.PaginationResponse{
				PageSize:   req.Pagination.PageSize,
				PageNumber: req.Pagination.PageNumber,
				Total:      0,
			},
			Comments: nil,
		}, nil
	}

	commenterIDs, commentableIDs, parentCommentIDs := s.gatherIDs(listResult.Comments)

	commenterMap, commentableMap, parentCommentMap, fetchErr := s.fetchEnrichmentData(ctx, commenterIDs, commentableIDs, parentCommentIDs)
	if fetchErr != nil {
		return commentparam.GetAllCommentsResponse{}, richerror.New(op).WithWarpError(fetchErr)
	}

	commentsInfo := s.buildDTOs(listResult.Comments, commenterMap, commentableMap, parentCommentMap)

	return commentparam.GetAllCommentsResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      listResult.TotalCount,
		},
		Comments: commentsInfo,
	}, nil

}

func (s Service) gatherIDs(comments []entity.Comment) (map[string][]string, map[string][]string, []string) {
	commenterIDs := make(map[string][]string)
	commentableIDs := make(map[string][]string)
	parentCommentIDs := make([]string, 0)

	for _, c := range comments {
		if c.CommenterID != "" {
			commenterIDs[c.CommenterType] = append(commenterIDs[c.CommenterType], c.CommenterID)
		}
		if c.CommentableID != "" {
			commentableIDs[c.CommentableType] = append(commentableIDs[c.CommentableType], c.CommentableID)
		}
		if c.ParentID != "" {
			parentCommentIDs = append(parentCommentIDs, c.ParentID)
		}
	}
	return commenterIDs, commentableIDs, parentCommentIDs
}

func (s Service) fetchEnrichmentData(ctx context.Context, commenterIDs, commentableIDs map[string][]string, parentCommentIDs []string) (map[string]Commenter, map[string]Commentable, map[string]entity.Comment, error) {
	var commenterMap map[string]Commenter
	var commentableMap map[string]Commentable
	var parentCommentMap map[string]entity.Comment

	errChan := make(chan error, 3)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		var err error
		commenterMap, err = s.commenterFetcher.FindCommentersByIDs(ctx, commenterIDs)
		if err != nil {
			errChan <- err
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		commentableMap, err = s.commentableFetcher.FindCommentablesByIDs(ctx, commentableIDs)
		if err != nil {
			errChan <- err
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		parentCommentMap, err = s.repo.FindManyByIDs(ctx, parentCommentIDs)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for fetchErr := range errChan {
		if fetchErr != nil {
			return nil, nil, nil, fetchErr
		}
	}
	return commenterMap, commentableMap, parentCommentMap, nil
}

func (s Service) buildDTOs(comments []entity.Comment, commenterMap map[string]Commenter, commentableMap map[string]Commentable, parentCommentMap map[string]entity.Comment) []commentparam.CommentsInfo {
	commentsInfo := make([]commentparam.CommentsInfo, 0, len(comments))
	for _, c := range comments {
		info := commentparam.CommentsInfo{
			ID:         c.ID,
			Body:       c.Body,
			IsApproved: c.IsApproved,
			CreatedAt:  c.CreatedAt,
		}
		if commenter, ok := commenterMap[c.CommenterID]; ok {
			info.Commenter.ID = commenter.GetID()
			info.Commenter.Name = commenter.GetName()
			info.Commenter.Type = c.CommenterType
		}
		if commentable, ok := commentableMap[c.CommentableID]; ok {
			info.Commentable.ID = commentable.GetID()
			info.Commentable.Title = commentable.GetTitle()
			info.Commentable.Type = c.CommentableType
		}
		if parent, ok := parentCommentMap[c.ParentID]; ok {
			info.ParentComment = &struct {
				ID            string `json:"id"`
				Body          string `json:"body"`
				CommenterName string `json:"commenter_name"`
			}{
				ID:            parent.ID,
				Body:          parent.Body,
				CommenterName: parent.CommenterName,
			}
		}
		commentsInfo = append(commentsInfo, info)
	}
	return commentsInfo
}

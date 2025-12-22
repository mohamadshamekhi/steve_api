package commentservice

import (
	"context"

	commentparam "github.com/abolfazlnorzad/steve/param/comment"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StorefrontService struct {
	repo               Repository
	commentableFetcher CommentableFetcher
	commenterFetcher   CommenterFetcher
}

func NewStorefrontService(
	repo Repository,
	commentableFetcher CommentableFetcher,
	commenterFetcher CommenterFetcher,
) StorefrontService {
	return StorefrontService{
		repo:               repo,
		commentableFetcher: commentableFetcher,
		commenterFetcher:   commenterFetcher,
	}
}

func (s StorefrontService) GetAprovedCommentsWithChildren(ctx context.Context, req commentparam.StorefrontGetCommentsRequest) (commentparam.StorefrontGetCommentsResponse, error) {
	const op = "commentService.GetAprovedCommentsWithChildren"

	lr, t, err := s.repo.GetApprovedCommentsWithChildren(ctx, req.CommentableType, req.CommentableID, req.Pagination)
	if err != nil {
		return commentparam.StorefrontGetCommentsResponse{}, richerror.New(op).WithWarpError(err)
	}

	return commentparam.StorefrontGetCommentsResponse{
		Comments:                 lr.Comments,
		AllApprovedCommentsCount: t,
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      lr.TotalCount,
		},
	}, nil
}

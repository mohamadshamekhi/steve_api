package ratingservice

import (
	"context"

	ratingparam "github.com/abolfazlnorzad/steve/param/rating"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type Service struct {
	repo Repository
}
type Repository interface {
	Upsert(ctx context.Context, UserID string, ProductID string, Rating float64) error
}

func New(r Repository) Service {
	return Service{
		repo: r,
	}
}

func (s Service) UpsertRating(ctx context.Context, req ratingparam.UpsertRatingRequest) (ratingparam.UpsertRatingResponse, error) {
	const op = "ratingservice.UpsertRating"
	if err := s.repo.Upsert(ctx, req.UserID, req.ProductID, req.Rating); err != nil {
		return ratingparam.UpsertRatingResponse{}, richerror.New(op).WithWarpError(err)
	}
	return ratingparam.UpsertRatingResponse{}, nil
}

package ratingservice

import (
	"context"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	ratingparam "github.com/abolfazlnorzad/steve/param/rating"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type Service struct {
	repo            Repository
	productProvider ProductProvider
}
type ProductProvider interface {
	UpdateProductRating(ctx context.Context, req productparam.UpdateProductRatingRequest) (productparam.UpdateProductRatingResponse, error)
}
type Repository interface {
	Upsert(ctx context.Context, UserID string, ProductID string, Rating float64) (float64, error)
}

func New(r Repository, productProvider ProductProvider) Service {
	return Service{
		repo:            r,
		productProvider: productProvider,
	}
}

func (s Service) UpsertRating(ctx context.Context, req ratingparam.UpsertRatingRequest) (ratingparam.UpsertRatingResponse, error) {
	const op = "ratingservice.UpsertRating"
	averageRating, err := s.repo.Upsert(ctx, req.UserID, req.ProductID, req.Rating)
	if err != nil {
		return ratingparam.UpsertRatingResponse{}, richerror.New(op).WithWarpError(err)
	}
	_, err = s.productProvider.UpdateProductRating(ctx, productparam.UpdateProductRatingRequest{ProductID: req.ProductID, AverageRating: averageRating})
	if err != nil {
		return ratingparam.UpsertRatingResponse{}, richerror.New(op).WithWarpError(err)
	}
	return ratingparam.UpsertRatingResponse{
		AverageRating: averageRating,
	}, nil
}

package mongodbcart

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	cartservice "github.com/abolfazlnorzad/steve/service/cart"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) GetAll(ctx context.Context, f cartparam.BackofficeFilterCart, p commonparam.PaginationRequest) (cartservice.ListResult, error) {
	const op = "mongodbcart.GetAll"
	filter := bson.M{}
	if f.UserID != nil {
		filter["user_id"] = *f.UserID
	}
	if f.Status != nil {
		filter["status"] = *f.Status
	}

	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return cartservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	offset := p.GetOffset()
	limit := p.GetPageSize()
	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)

	if err != nil {
		return cartservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var carts []entity.Cart
	if err := cursor.All(ctx, &carts); err != nil {
		return cartservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return cartservice.ListResult{
		Carts:      carts,
		TotalCount: totalCount,
	}, nil
}

package mongodborder

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	orderservice "github.com/abolfazlnorzad/steve/service/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) GetAll(ctx context.Context, f orderparam.BackofficeFilterOrder, p commonparam.PaginationRequest) (orderservice.ListResult, error) {
	const op = "mongodborder.GetAll"
	filter := bson.M{}
	if f.UserID != nil {
		filter["user_id"] = *f.UserID
	}
	if f.Status != nil {
		filter["status"] = *f.Status
	}

	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return orderservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	offset := p.GetOffset()
	limit := p.GetPageSize()
	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)

	if err != nil {
		return orderservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var orders []entity.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return orderservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return orderservice.ListResult{
		Orders:     orders,
		TotalCount: totalCount,
	}, nil
}

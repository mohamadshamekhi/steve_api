package mongodborder

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	orderservice "github.com/abolfazlnorzad/steve/service/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AllPartnerOrders lists orders placed from the partner shop, newest first,
// optionally filtered by status.
func (d DB) AllPartnerOrders(ctx context.Context, status string, p commonparam.PaginationRequest) (orderservice.ListResult, error) {
	const op = "mongodborder.AllPartnerOrders"
	filter := bson.M{"is_partner_order": true}
	if status != "" {
		filter["status"] = status
	}

	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return orderservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	findOptions := options.Find().
		SetSkip(int64(p.GetOffset())).
		SetLimit(int64(p.GetPageSize())).
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

	return orderservice.ListResult{Orders: orders, TotalCount: totalCount}, nil
}

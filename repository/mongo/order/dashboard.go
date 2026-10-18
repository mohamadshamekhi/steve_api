package mongodborder

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	dashboardparam "github.com/abolfazlnorzad/steve/param/dashboard"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createdBetween(from, to time.Time) bson.M {
	return bson.M{"created_at": bson.M{"$gte": from, "$lte": to}}
}

// CountByStatus counts the orders registered in the window, per status.
func (d DB) CountByStatus(ctx context.Context, from, to time.Time) ([]dashboardparam.OrderStatusCount, error) {
	const op = "mongodborder.CountByStatus"

	pipeline := []bson.M{
		{"$match": createdBetween(from, to)},
		{"$group": bson.M{"_id": "$status", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}},
	}
	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	rows := make([]dashboardparam.OrderStatusCount, 0)
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return rows, nil
}

// Recent returns the newest orders, newest first.
func (d DB) Recent(ctx context.Context, limit int64) ([]entity.Order, error) {
	const op = "mongodborder.Recent"

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetProjection(bson.M{"items": 0, "shipping_address": 0})
	cursor, err := d.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	orders := make([]entity.Order, 0)
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return orders, nil
}

// CancelledOrdersBetween lists the orders cancelled in the window (matched
// by created_at, same as the rest of this file). Used to break down
// cancellations by reason/payment status for reporting.
func (d DB) CancelledOrdersBetween(ctx context.Context, from, to time.Time) ([]entity.Order, error) {
	const op = "mongodborder.CancelledOrdersBetween"

	filter := createdBetween(from, to)
	filter["status"] = entity.OrderStatusCancelled

	cursor, err := d.collection.Find(ctx, filter, options.Find().SetProjection(bson.M{"items": 0, "shipping_address": 0}))
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	orders := make([]entity.Order, 0)
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return orders, nil
}

// SoldOrdersBetween lists the orders in the window that count as a sale
// (neither awaiting payment nor cancelled). Used to backfill the sales
// ledger for orders registered before it existed.
func (d DB) SoldOrdersBetween(ctx context.Context, from, to time.Time) ([]entity.Order, error) {
	const op = "mongodborder.SoldOrdersBetween"

	filter := createdBetween(from, to)
	filter["status"] = bson.M{"$nin": bson.A{entity.OrderStatusPendingPayment, entity.OrderStatusCancelled}}

	cursor, err := d.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	orders := make([]entity.Order, 0)
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return orders, nil
}

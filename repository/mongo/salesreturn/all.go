package mongodbsalesreturn

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	salesreturnparam "github.com/abolfazlnorzad/steve/param/salesreturn"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// exportLimit caps unpaginated reads (CSV export).
const exportLimit = 20000

func buildFilter(f salesreturnparam.SalesReturnFilter) bson.M {
	filter := bson.M{}
	if f.OrderID != "" {
		filter["order_id"] = f.OrderID
	}
	if f.UserID != "" {
		filter["user_id"] = f.UserID
	}
	if f.From != nil || f.To != nil {
		returnedAt := bson.M{}
		if f.From != nil {
			returnedAt["$gte"] = *f.From
		}
		if f.To != nil {
			returnedAt["$lte"] = *f.To
		}
		filter["returned_at"] = returnedAt
	}
	return filter
}

func (d DB) All(ctx context.Context, f salesreturnparam.SalesReturnFilter, p commonparam.PaginationRequest) ([]entity.SalesReturn, int64, error) {
	const op = "mongodbsalesreturn.All"

	filter := buildFilter(f)

	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	findOptions := options.Find().
		SetSkip(int64(p.GetOffset())).
		SetLimit(int64(p.GetPageSize())).
		SetSort(bson.D{{Key: "returned_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	returns := make([]entity.SalesReturn, 0)
	if err := cursor.All(ctx, &returns); err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return returns, total, nil
}

func (d DB) AllForExport(ctx context.Context, f salesreturnparam.SalesReturnFilter) ([]entity.SalesReturn, error) {
	const op = "mongodbsalesreturn.AllForExport"

	findOptions := options.Find().
		SetLimit(exportLimit).
		SetSort(bson.D{{Key: "returned_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, buildFilter(f), findOptions)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	returns := make([]entity.SalesReturn, 0)
	if err := cursor.All(ctx, &returns); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return returns, nil
}

// Summary totals the returns in the window: the «برگشت از فروش» line of the
// profit and loss statement.
func (d DB) Summary(ctx context.Context, f salesreturnparam.SalesReturnFilter) (salesreturnparam.SalesReturnSummary, error) {
	const op = "mongodbsalesreturn.Summary"

	pipeline := []bson.M{
		{"$match": buildFilter(f)},
		{"$group": bson.M{
			"_id":           nil,
			"count":         bson.M{"$sum": 1},
			"refund_amount": bson.M{"$sum": "$refund_amount"},
			"items_amount":  bson.M{"$sum": "$items_amount"},
			"returned_cost": bson.M{"$sum": "$returned_cost"},
		}},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return salesreturnparam.SalesReturnSummary{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var rows []struct {
		Count        int64  `bson:"count"`
		RefundAmount uint64 `bson:"refund_amount"`
		ItemsAmount  uint64 `bson:"items_amount"`
		ReturnedCost uint64 `bson:"returned_cost"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return salesreturnparam.SalesReturnSummary{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if len(rows) == 0 {
		return salesreturnparam.SalesReturnSummary{}, nil
	}
	return salesreturnparam.SalesReturnSummary{
		Count:        rows[0].Count,
		RefundAmount: rows[0].RefundAmount,
		ItemsAmount:  rows[0].ItemsAmount,
		ReturnedCost: rows[0].ReturnedCost,
	}, nil
}

package mongodbcash

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func buildFilter(f accountingparam.CashFilter) bson.M {
	filter := bson.M{}
	if f.Direction != "" {
		filter["direction"] = f.Direction
	}
	if f.Method != "" {
		filter["method"] = f.Method
	}
	if f.RefType != "" {
		filter["ref_type"] = f.RefType
	}
	if f.RefID != "" {
		filter["ref_id"] = f.RefID
	}
	if f.SupplierID != "" {
		filter["supplier_id"] = f.SupplierID
	}
	if f.From != nil || f.To != nil {
		occurredAt := bson.M{}
		if f.From != nil {
			occurredAt["$gte"] = *f.From
		}
		if f.To != nil {
			occurredAt["$lte"] = *f.To
		}
		filter["occurred_at"] = occurredAt
	}
	return filter
}

func (d DB) All(ctx context.Context, f accountingparam.CashFilter, p commonparam.PaginationRequest) ([]entity.CashTransaction, int64, error) {
	const op = "mongodbcash.All"

	filter := buildFilter(f)

	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	findOptions := options.Find().
		SetSkip(int64(p.GetOffset())).
		SetLimit(int64(p.GetPageSize())).
		SetSort(bson.D{{Key: "occurred_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	transactions := make([]entity.CashTransaction, 0)
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return transactions, total, nil
}

func (d DB) AllForExport(ctx context.Context, f accountingparam.CashFilter) ([]entity.CashTransaction, error) {
	const op = "mongodbcash.AllForExport"

	findOptions := options.Find().
		SetLimit(exportLimit).
		SetSort(bson.D{{Key: "occurred_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, buildFilter(f), findOptions)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	transactions := make([]entity.CashTransaction, 0)
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return transactions, nil
}

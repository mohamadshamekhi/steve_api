package mongodbexpense

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// buildFilter turns the API filter into a Mongo query. It is shared by the
// list and the aggregations so a report can never disagree with the list it
// was opened from.
func buildFilter(f accountingparam.ExpenseFilter) bson.M {
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}
	if f.Category != "" {
		filter["category"] = f.Category
	}
	if f.SupplierID != "" {
		filter["supplier_id"] = f.SupplierID
	}
	if f.IsPaid != nil {
		filter["is_paid"] = *f.IsPaid
	}
	if f.Search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": f.Search, "$options": "i"}},
			{"note": bson.M{"$regex": f.Search, "$options": "i"}},
		}
	}
	if f.From != nil || f.To != nil {
		spentAt := bson.M{}
		if f.From != nil {
			spentAt["$gte"] = *f.From
		}
		if f.To != nil {
			spentAt["$lte"] = *f.To
		}
		filter["spent_at"] = spentAt
	}
	return filter
}

func (d DB) All(ctx context.Context, f accountingparam.ExpenseFilter, p commonparam.PaginationRequest) ([]entity.Expense, int64, error) {
	const op = "mongodbexpense.All"

	filter := buildFilter(f)

	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	findOptions := options.Find().
		SetSkip(int64(p.GetOffset())).
		SetLimit(int64(p.GetPageSize())).
		SetSort(bson.D{{Key: "spent_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	expenses := make([]entity.Expense, 0)
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return expenses, total, nil
}

// AllForExport returns every matching row without pagination, capped so a
// careless date range cannot pull the whole collection into memory.
func (d DB) AllForExport(ctx context.Context, f accountingparam.ExpenseFilter) ([]entity.Expense, error) {
	const op = "mongodbexpense.AllForExport"

	findOptions := options.Find().
		SetLimit(exportLimit).
		SetSort(bson.D{{Key: "spent_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, buildFilter(f), findOptions)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	expenses := make([]entity.Expense, 0)
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return expenses, nil
}

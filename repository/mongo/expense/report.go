package mongodbexpense

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// exportLimit caps unpaginated reads (CSV export) at a size that still fits
// comfortably in memory and in a spreadsheet.
const exportLimit = 20000

// Total sums the matching expenses — the «هزینه‌های عملیاتی» line of the
// profit and loss statement.
func (d DB) Total(ctx context.Context, f accountingparam.ExpenseFilter) (uint64, error) {
	const op = "mongodbexpense.Total"

	pipeline := []bson.M{
		{"$match": buildFilter(f)},
		{"$group": bson.M{"_id": nil, "amount": bson.M{"$sum": "$amount"}}},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var rows []struct {
		Amount uint64 `bson:"amount"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if len(rows) == 0 {
		return 0, nil
	}
	return rows[0].Amount, nil
}

// TotalByCategory answers «پولم کجا رفت».
func (d DB) TotalByCategory(ctx context.Context, f accountingparam.ExpenseFilter) ([]accountingparam.ExpenseCategoryRow, error) {
	const op = "mongodbexpense.TotalByCategory"

	pipeline := []bson.M{
		{"$match": buildFilter(f)},
		{"$group": bson.M{
			"_id":    "$category",
			"amount": bson.M{"$sum": "$amount"},
			"count":  bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"amount": -1}},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var aggRows []struct {
		ID     string `bson:"_id"`
		Amount uint64 `bson:"amount"`
		Count  int64  `bson:"count"`
	}
	if err := cursor.All(ctx, &aggRows); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	rows := make([]accountingparam.ExpenseCategoryRow, 0, len(aggRows))
	for _, r := range aggRows {
		rows = append(rows, accountingparam.ExpenseCategoryRow{
			Category: entity.ExpenseCategory(r.ID),
			Amount:   r.Amount,
			Count:    r.Count,
		})
	}
	return rows, nil
}

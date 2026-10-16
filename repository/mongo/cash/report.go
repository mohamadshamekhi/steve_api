package mongodbcash

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// directionAccumulators split each bucket's total into money in and money
// out, so one pass over the collection answers both.
var directionAccumulators = bson.M{
	"total_in": bson.M{"$sum": bson.M{"$cond": []any{
		bson.M{"$eq": []any{"$direction", entity.CashIn}}, "$amount", 0,
	}}},
	"total_out": bson.M{"$sum": bson.M{"$cond": []any{
		bson.M{"$eq": []any{"$direction", entity.CashOut}}, "$amount", 0,
	}}},
}

type cashAggRow struct {
	ID       string `bson:"_id"`
	TotalIn  uint64 `bson:"total_in"`
	TotalOut uint64 `bson:"total_out"`
}

func (d DB) aggregate(ctx context.Context, op string, f accountingparam.CashFilter, groupID any) ([]cashAggRow, error) {
	group := bson.M{"_id": groupID}
	for k, v := range directionAccumulators {
		group[k] = v
	}

	pipeline := []bson.M{
		{"$match": buildFilter(f)},
		{"$group": group},
		{"$sort": bson.M{"_id": 1}},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var rows []cashAggRow
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return rows, nil
}

// Summary is the دخل و خرج total of the window.
func (d DB) Summary(ctx context.Context, f accountingparam.CashFilter) (accountingparam.CashSummary, error) {
	rows, err := d.aggregate(ctx, "mongodbcash.Summary", f, nil)
	if err != nil {
		return accountingparam.CashSummary{}, err
	}
	if len(rows) == 0 {
		return accountingparam.CashSummary{}, nil
	}
	return accountingparam.CashSummary{
		TotalIn:  rows[0].TotalIn,
		TotalOut: rows[0].TotalOut,
		Balance:  int64(rows[0].TotalIn) - int64(rows[0].TotalOut),
	}, nil
}

// SummaryByMethod splits the window by how the money moved (نقد/کارت/درگاه).
func (d DB) SummaryByMethod(ctx context.Context, f accountingparam.CashFilter) ([]accountingparam.CashMethodRow, error) {
	aggRows, err := d.aggregate(ctx, "mongodbcash.SummaryByMethod", f, "$method")
	if err != nil {
		return nil, err
	}

	rows := make([]accountingparam.CashMethodRow, 0, len(aggRows))
	for _, r := range aggRows {
		rows = append(rows, accountingparam.CashMethodRow{
			Method:   entity.PaymentMethod(r.ID),
			TotalIn:  r.TotalIn,
			TotalOut: r.TotalOut,
			Balance:  int64(r.TotalIn) - int64(r.TotalOut),
		})
	}
	return rows, nil
}

// SummaryByDay is the daily till report (گردش روزانه صندوق).
func (d DB) SummaryByDay(ctx context.Context, f accountingparam.CashFilter) ([]accountingparam.CashDayRow, error) {
	dateToString := bson.M{"format": "%Y-%m-%d", "date": "$occurred_at"}
	if f.Timezone != "" {
		dateToString["timezone"] = f.Timezone
	}
	groupID := bson.M{"$dateToString": dateToString}

	aggRows, err := d.aggregate(ctx, "mongodbcash.SummaryByDay", f, groupID)
	if err != nil {
		return nil, err
	}

	rows := make([]accountingparam.CashDayRow, 0, len(aggRows))
	for _, r := range aggRows {
		rows = append(rows, accountingparam.CashDayRow{
			Day:      r.ID,
			TotalIn:  r.TotalIn,
			TotalOut: r.TotalOut,
			Balance:  int64(r.TotalIn) - int64(r.TotalOut),
		})
	}
	return rows, nil
}

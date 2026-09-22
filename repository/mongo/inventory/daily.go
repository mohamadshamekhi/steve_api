package mongodbinventory

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// DailySales buckets the sales ledger by calendar day in the given timezone
// (a Mongo offset such as "+03:30"). Days without sales are absent; the
// caller fills the gaps. OrderCount is the number of distinct orders.
func (d DB) DailySales(ctx context.Context, req inventoryparam.ProfitReportRequest, timezone string) ([]inventoryparam.DailySalesRow, error) {
	const op = "mongodbinventory.DailySales"

	pipeline := append(salesStages(req),
		bson.M{"$group": bson.M{
			"_id": bson.M{"$dateToString": bson.M{
				"format": "%Y-%m-%d", "date": "$created_at", "timezone": timezone,
			}},
			"revenue":      bson.M{"$sum": "$total_price"},
			"gross_profit": bson.M{"$sum": "$profit"},
			// Every sale line is written per order (ref_type=order).
			"orders": bson.M{"$addToSet": "$ref_id"},
		}},
		bson.M{"$project": bson.M{
			"revenue":      1,
			"gross_profit": 1,
			"order_count":  bson.M{"$size": "$orders"},
		}},
		bson.M{"$sort": bson.M{"_id": 1}},
	)

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	rows := make([]inventoryparam.DailySalesRow, 0)
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return rows, nil
}

// OrderIDsWithSales returns which of the given order ids already have sale
// lines in the ledger, so a backfill never books an order twice.
func (d DB) OrderIDsWithSales(ctx context.Context, orderIDs []string) (map[string]struct{}, error) {
	const op = "mongodbinventory.OrderIDsWithSales"

	result := make(map[string]struct{})
	if len(orderIDs) == 0 {
		return result, nil
	}

	ids, err := d.collection.Distinct(ctx, "ref_id", bson.M{
		"ref_type": entity.StockRefOrder,
		"ref_id":   bson.M{"$in": orderIDs},
		"type":     bson.M{"$in": saleTypes},
	})
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	for _, id := range ids {
		if s, ok := id.(string); ok {
			result[s] = struct{}{}
		}
	}
	return result, nil
}

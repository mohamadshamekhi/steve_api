package mongodbinventory

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// saleTypes are the movement types that represent revenue: a sale out of
// stock, a sale sourced externally, and a service sold.
var saleTypes = bson.A{
	entity.StockMovementSaleOut,
	entity.StockMovementExternalSupply,
	entity.StockMovementServiceSale,
}

// profitMatch builds the $match stage shared by the profit aggregations:
// sale movements only, inside the requested date window.
func profitMatch(req inventoryparam.ProfitReportRequest) bson.M {
	match := bson.M{"type": bson.M{"$in": saleTypes}}
	if req.From != nil || req.To != nil {
		createdAt := bson.M{}
		if req.From != nil {
			createdAt["$gte"] = *req.From
		}
		if req.To != nil {
			createdAt["$lte"] = *req.To
		}
		match["created_at"] = createdAt
	}
	return match
}

// profitAccumulators are the summed fields every profit bucket needs.
// Quantity is stored negative for sales, so units are negated back.
var profitAccumulators = bson.M{
	"units_sold":   bson.M{"$sum": bson.M{"$abs": "$quantity"}},
	"revenue":      bson.M{"$sum": "$total_price"},
	"cost":         bson.M{"$sum": "$total_cost"},
	"gross_profit": bson.M{"$sum": "$profit"},
}

type profitAggRow struct {
	ID          any    `bson:"_id"`
	Label       string `bson:"label"`
	UnitsSold   int64  `bson:"units_sold"`
	Revenue     uint64 `bson:"revenue"`
	Cost        uint64 `bson:"cost"`
	GrossProfit int64  `bson:"gross_profit"`
}

// ProfitSummary returns the single totals row for the window.
func (d DB) ProfitSummary(ctx context.Context, req inventoryparam.ProfitReportRequest) (inventoryparam.ProfitRow, error) {
	const op = "mongodbinventory.ProfitSummary"

	group := bson.M{"_id": nil}
	for k, v := range profitAccumulators {
		group[k] = v
	}

	pipeline := []bson.M{
		{"$match": profitMatch(req)},
		{"$group": group},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return inventoryparam.ProfitRow{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var rows []profitAggRow
	if err := cursor.All(ctx, &rows); err != nil {
		return inventoryparam.ProfitRow{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if len(rows) == 0 {
		return inventoryparam.ProfitRow{}, nil
	}
	return inventoryparam.ProfitRow{
		UnitsSold:   rows[0].UnitsSold,
		Revenue:     rows[0].Revenue,
		Cost:        rows[0].Cost,
		GrossProfit: rows[0].GrossProfit,
	}, nil
}

// ProfitBreakdown groups the same window by day, product or supplier.
func (d DB) ProfitBreakdown(ctx context.Context, req inventoryparam.ProfitReportRequest) ([]inventoryparam.ProfitRow, error) {
	const op = "mongodbinventory.ProfitBreakdown"

	var groupID any
	label := bson.M{"$first": ""}

	switch req.GroupBy {
	case "day":
		groupID = bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$created_at"}}
	case "supplier":
		groupID = "$supplier_id"
	default: // product
		groupID = "$product_id"
		label = bson.M{"$first": "$product_name"}
	}

	group := bson.M{"_id": groupID, "label": label}
	for k, v := range profitAccumulators {
		group[k] = v
	}

	pipeline := []bson.M{
		{"$match": profitMatch(req)},
		{"$group": group},
		{"$sort": bson.M{"gross_profit": -1}},
		{"$limit": 500},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var aggRows []profitAggRow
	if err := cursor.All(ctx, &aggRows); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	rows := make([]inventoryparam.ProfitRow, 0, len(aggRows))
	for _, r := range aggRows {
		key, _ := r.ID.(string)
		rows = append(rows, inventoryparam.ProfitRow{
			Key:         key,
			Label:       r.Label,
			UnitsSold:   r.UnitsSold,
			Revenue:     r.Revenue,
			Cost:        r.Cost,
			GrossProfit: r.GrossProfit,
		})
	}
	return rows, nil
}

package mongodbproduct

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// variantArrayFilter targets one variant inside the product's variants array.
func variantArrayFilter(variantID string) *options.UpdateOptions {
	return options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []any{bson.M{"elem.id": variantID}},
	})
}

// SetVariantStockConfig writes the inventory/sourcing settings of one variant.
// Only the non-nil fields of req are touched.
func (d DB) SetVariantStockConfig(ctx context.Context, req inventoryparam.SetStockConfigRequest) error {
	const op = "mongodbproduct.SetVariantStockConfig"

	set := bson.M{"updated_at": time.Now()}
	unset := bson.M{}

	if req.StockMode != nil {
		set["variants.$[elem].stock_mode"] = *req.StockMode
	}
	if req.AllowBackorder != nil {
		set["variants.$[elem].allow_backorder"] = *req.AllowBackorder
	}
	if req.BackorderLimit != nil {
		set["variants.$[elem].backorder_limit"] = *req.BackorderLimit
	}
	if req.SupplierID != nil {
		if *req.SupplierID == "" {
			unset["variants.$[elem].supplier_id"] = ""
		} else {
			set["variants.$[elem].supplier_id"] = *req.SupplierID
		}
	}
	if req.SupplyCost != nil {
		set["variants.$[elem].supply_cost"] = *req.SupplyCost
	}
	if req.AvgCost != nil {
		set["variants.$[elem].avg_cost"] = *req.AvgCost
	}

	update := bson.M{"$set": set}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	result, err := d.collection.UpdateOne(ctx, bson.M{"_id": req.ProductID}, update, variantArrayFilter(req.VariantID))
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("product.not_found").
			WithMeta(map[string]any{"product_id": req.ProductID})
	}
	return nil
}

// ApplyStockDelta moves a variant's quantity by a signed delta and, when
// avgCost is non-nil, writes the recalculated weighted-average cost in the
// same update — the two must never drift apart.
//
// Quantity is a uint64 in the entity but Mongo's $inc happily goes negative;
// callers must not push it below zero except for a deliberate backorder.
func (d DB) ApplyStockDelta(ctx context.Context, productID string, variantID string, delta int64, avgCost *uint64) error {
	const op = "mongodbproduct.ApplyStockDelta"

	update := bson.M{
		"$inc": bson.M{"variants.$[elem].quantity": delta},
		"$set": bson.M{"updated_at": time.Now()},
	}
	if avgCost != nil {
		update["$set"].(bson.M)["variants.$[elem].avg_cost"] = *avgCost
	}

	result, err := d.collection.UpdateOne(ctx, bson.M{"_id": productID}, update, variantArrayFilter(variantID))
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("product.not_found").
			WithMeta(map[string]any{"product_id": productID, "variant_id": variantID})
	}
	if result.ModifiedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("product.variant_not_found").
			WithMeta(map[string]any{"product_id": productID, "variant_id": variantID})
	}
	return nil
}

// InventoryValue lists every stock-tracked variant with its current valuation,
// newest products first, and returns the whole catalog's totals alongside the
// requested page.
func (d DB) InventoryValue(ctx context.Context, req inventoryparam.InventoryValueRequest) (inventoryparam.InventoryValueResponse, error) {
	const op = "mongodbproduct.InventoryValue"

	// Only physical products hold stock; services never do.
	match := bson.M{
		"deleted_at": bson.M{"$exists": false},
		"kind":       bson.M{"$ne": entity.ProductKindService},
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$unwind": "$variants"},
		// stock_mode is absent on pre-inventory documents, which were tracked.
		{"$match": bson.M{"variants.stock_mode": bson.M{"$ne": entity.StockModeOnDemand}}},
	}

	if req.LowStockOnly {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{"variants.quantity": bson.M{"$lte": req.LowStockThreshold}},
		})
	}

	pipeline = append(pipeline, bson.M{"$facet": bson.M{
		"rows": []bson.M{
			{"$sort": bson.M{"created_at": -1}},
			{"$skip": int64(req.Pagination.GetOffset())},
			{"$limit": int64(req.Pagination.GetPageSize())},
			{"$project": bson.M{
				"_id":          0,
				"product_id":   "$_id",
				"product_name": "$name",
				"variant_id":   "$variants.id",
				"attributes":   "$variants.variant_attribute",
				"stock_mode":   "$variants.stock_mode",
				"quantity":     "$variants.quantity",
				"avg_cost":     "$variants.avg_cost",
				"price":        "$variants.price",
			}},
		},
		"totals": []bson.M{
			{"$group": bson.M{
				"_id":         nil,
				"total_units": bson.M{"$sum": "$variants.quantity"},
				"total_value": bson.M{"$sum": bson.M{"$multiply": bson.A{"$variants.quantity", bson.M{"$ifNull": bson.A{"$variants.avg_cost", 0}}}}},
				"count":       bson.M{"$sum": 1},
			}},
		},
	}})

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return inventoryparam.InventoryValueResponse{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var facet []struct {
		Rows   []inventoryparam.InventoryValueRow `bson:"rows"`
		Totals []struct {
			TotalUnits uint64 `bson:"total_units"`
			TotalValue uint64 `bson:"total_value"`
			Count      int64  `bson:"count"`
		} `bson:"totals"`
	}
	if err := cursor.All(ctx, &facet); err != nil {
		return inventoryparam.InventoryValueResponse{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	resp := inventoryparam.InventoryValueResponse{
		Rows: make([]inventoryparam.InventoryValueRow, 0),
		Pagination: commonparam.PaginationResponse{
			PageNumber: req.Pagination.GetPageNumber(),
			PageSize:   req.Pagination.GetPageSize(),
		},
	}
	if len(facet) == 0 {
		return resp, nil
	}

	resp.Rows = facet[0].Rows
	for i := range resp.Rows {
		row := &resp.Rows[i]
		if row.StockMode == "" {
			row.StockMode = entity.StockModeTracked
		}
		row.StockValue = row.Quantity * row.AvgCost
		row.PotentialProfit = (int64(row.Price) - int64(row.AvgCost)) * int64(row.Quantity)
	}
	if len(facet[0].Totals) > 0 {
		resp.TotalUnits = facet[0].Totals[0].TotalUnits
		resp.TotalStockValue = facet[0].Totals[0].TotalValue
		resp.Pagination.Total = uint(facet[0].Totals[0].Count)
	}
	return resp, nil
}

package mongodbinventory

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func buildMovementFilter(f inventoryparam.MovementFilter) bson.M {
	filter := bson.M{}
	if f.ProductID != "" {
		filter["product_id"] = f.ProductID
	}
	if f.VariantID != "" {
		filter["variant_id"] = f.VariantID
	}
	if f.SupplierID != "" {
		filter["supplier_id"] = f.SupplierID
	}
	if f.Type != "" {
		filter["type"] = f.Type
	}
	if f.RefType != "" {
		filter["ref_type"] = f.RefType
	}
	if f.RefID != "" {
		filter["ref_id"] = f.RefID
	}
	if f.From != nil || f.To != nil {
		createdAt := bson.M{}
		if f.From != nil {
			createdAt["$gte"] = *f.From
		}
		if f.To != nil {
			createdAt["$lte"] = *f.To
		}
		filter["created_at"] = createdAt
	}
	return filter
}

// All lists ledger movements newest-first.
func (d DB) All(ctx context.Context, f inventoryparam.MovementFilter, p commonparam.PaginationRequest) ([]entity.StockMovement, int64, error) {
	const op = "mongodbinventory.All"

	filter := buildMovementFilter(f)

	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	findOptions := options.Find().
		SetSkip(int64(p.GetOffset())).
		SetLimit(int64(p.GetPageSize())).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	movements := make([]entity.StockMovement, 0)
	if err := cursor.All(ctx, &movements); err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return movements, total, nil
}

// DeleteByRef removes every movement written for one reference. Used when a
// purchase invoice is cancelled and its ledger effect is reversed by a
// compensating write, so the original lines must not be double counted.
func (d DB) DeleteByRef(ctx context.Context, refType string, refID string) error {
	const op = "mongodbinventory.DeleteByRef"

	if _, err := d.collection.DeleteMany(ctx, bson.M{"ref_type": refType, "ref_id": refID}); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return nil
}

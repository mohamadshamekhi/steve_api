package mongodbcart

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpsertCustomItem writes a hand-typed invoice line into the user's active
// cart. When item.VariantID is already set the matching line is rewritten in
// place (that is what the panel's +/- buttons and inline edits use);
// otherwise a new line is pushed and its generated id returned, which is the
// only handle the caller has on a line that has no product behind it.
func (d DB) UpsertCustomItem(ctx context.Context, userID string, item entity.CartItem) (string, error) {
	const op = "mongodbcart.UpsertCustomItem"

	item.IsCustom = true
	item.ProductID = ""
	if item.Quantity < 1 {
		item.Quantity = 1
	}
	if item.Kind == "" {
		item.Kind = entity.ProductKindPhysical
	}

	if item.VariantID != "" {
		filter := bson.M{
			"user_id":          userID,
			"status":           entity.CartActive,
			"items.variant_id": item.VariantID,
			"items.is_custom":  true,
		}
		update := bson.M{
			"$set": bson.M{
				"items.$.quantity":      item.Quantity,
				"items.$.price_at_time": item.PriceAtTime,
				"items.$.product_name":  item.ProductName,
				"items.$.kind":          item.Kind,
				"items.$.unit_cost":     item.UnitCost,
				"updated_at":            time.Now(),
			},
		}
		res, err := d.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return "", richerror.New(op).WithWarpError(err).
				WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
		}
		if res.MatchedCount == 0 {
			return "", richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("cart.item_not_found")
		}
		return item.VariantID, nil
	}

	item.VariantID = ulid.Make().String()
	filterAdd := bson.M{"user_id": userID, "status": entity.CartActive}
	updateAdd := bson.M{
		"$push": bson.M{"items": item},
		"$setOnInsert": bson.M{
			"_id":        ulid.Make().String(),
			"user_id":    userID,
			"status":     entity.CartActive,
			"created_at": time.Now(),
		},
		"$set": bson.M{"updated_at": time.Now()},
	}
	if _, err := d.collection.UpdateOne(ctx, filterAdd, updateAdd, options.Update().SetUpsert(true)); err != nil {
		return "", richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return item.VariantID, nil
}

// SetItemManualDiscount stores the counter discount given on one line of the
// user's active cart, in Toman for the whole line.
func (d DB) SetItemManualDiscount(ctx context.Context, userID string, variantID string, amount int64) error {
	const op = "mongodbcart.SetItemManualDiscount"

	filter := bson.M{
		"user_id":          userID,
		"status":           entity.CartActive,
		"items.variant_id": variantID,
	}
	update := bson.M{
		"$set": bson.M{
			"items.$.manual_discount": amount,
			"updated_at":              time.Now(),
		},
	}
	res, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if res.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("cart.item_not_found")
	}
	return nil
}

// SetCartManualDiscount stores the counter discount given on the invoice as a
// whole for the user's active cart.
func (d DB) SetCartManualDiscount(ctx context.Context, userID string, amount uint) error {
	const op = "mongodbcart.SetCartManualDiscount"

	filter := bson.M{"user_id": userID, "status": entity.CartActive}
	update := bson.M{"$set": bson.M{"manual_discount": amount, "updated_at": time.Now()}}

	res, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if res.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("cart.not_found")
	}
	return nil
}

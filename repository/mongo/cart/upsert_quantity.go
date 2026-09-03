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

// UpsertItemWithQuantity sets an item's quantity (and price) to an absolute
// value in the user's active cart, creating the cart/item if needed. Unlike
// AddToCart (which only ever +1s), this is what the partner shop needs so a
// colleague can order N units at once.
func (d DB) UpsertItemWithQuantity(ctx context.Context, userID string, item entity.CartItem, quantity int) error {
	const op = "mongodbcart.UpsertItemWithQuantity"
	if quantity < 1 {
		quantity = 1
	}

	setExisting := bson.M{
		"user_id":          userID,
		"status":           entity.CartActive,
		"items.product_id": item.ProductID,
		"items.variant_id": item.VariantID,
	}
	updateExisting := bson.M{
		"$set": bson.M{
			"items.$.quantity":      quantity,
			"items.$.price_at_time": item.PriceAtTime,
			"updated_at":            time.Now(),
		},
	}
	res, err := d.collection.UpdateOne(ctx, setExisting, updateExisting)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if res.MatchedCount > 0 {
		return nil
	}

	item.Quantity = quantity
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
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return nil
}

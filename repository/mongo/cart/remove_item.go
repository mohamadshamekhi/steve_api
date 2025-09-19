package mongodbcart

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) RemoveFromCart(ctx context.Context, userID string, productID string, variantID string) error {
	const op = "mongodbcart.RemoveFromCart"

	filterDecrement := bson.M{
		"user_id":          userID,
		"status":           "active",
		"items.product_id": productID,
		"items.variant_id": variantID,
	}
	updateDecrement := bson.M{
		"$inc": bson.M{"items.$.quantity": -1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	result, err := d.collection.UpdateOne(ctx, filterDecrement, updateDecrement)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if result.MatchedCount == 0 {
		return nil
	}

	filterRemove := bson.M{
		"user_id": userID,
		"status":  "active",
	}
	updateRemove := bson.M{
		"$pull": bson.M{"items": bson.M{"quantity": bson.M{"$lte": 0}}},
	}

	_, err = d.collection.UpdateOne(ctx, filterRemove, updateRemove)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	return nil
}

func (d DB) DeleteItemFromCart(ctx context.Context, userID string, variantID string) error {
	const op = "mongodbcart.DeleteItemFromCart"
	filter := bson.M{
		"user_id": userID,
		"status":  entity.CartActive,
	}

	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{"variant_id": variantID},
		},
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("cart.not_found")
	}

	return nil
}

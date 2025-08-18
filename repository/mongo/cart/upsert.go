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

func (d DB) AddToCart(ctx context.Context, userID string, item entity.CartItem) error {
	const op = "mongodbcart.AddToCart"
	filterIncrement := bson.M{
		"user_id":          userID,
		"status":           entity.CartActive,
		"items.product_id": item.ProductID,
		"items.variant_id": item.VariantID,
	}

	updateIncrement := bson.M{
		"$inc": bson.M{"items.$.quantity": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	result, err := d.collection.UpdateOne(ctx, filterIncrement, updateIncrement)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if result.MatchedCount > 0 {
		return nil
	}

	item.Quantity = 1

	filterAdd := bson.M{
		"user_id": userID,
		"status":  entity.CartActive,
	}

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

	opts := options.Update().SetUpsert(true)

	_, err = d.collection.UpdateOne(ctx, filterAdd, updateAdd, opts)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	return nil
}

package mongodbcart

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) UpdateStatus(ctx context.Context, cartID string, status entity.CartStatus) error {
	const op = "mongodbcart.UpdateStatus"
	filter := bson.M{"_id": cartID}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("cart.not_found").
			WithMeta(map[string]any{"cart_id": cartID})
	}

	return nil

}

func (d DB) UpdateItemPrice(ctx context.Context, userID string, variantID string, newPrice int64) error {
	const op = "mongodbcart.UpdateItemPrice"

	filter := bson.M{
		"user_id": userID,
		"status":  entity.CartActive,
	}

	update := bson.M{
		"$set": bson.M{"items.$[elem].price_at_time": newPrice},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem.variant_id": variantID}},
	})

	result, err := d.collection.UpdateOne(ctx, filter, update, arrayFilters)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("cart.item_not_found")
	}

	return nil
}

package mongodbrating

import (
	"context"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) Upsert(ctx context.Context, UserID string, ProductID string, Rating float64) error {
	const op = "mongodbrating.Upsert"

	id := ulid.Make().String()
	filter := bson.M{
		"user_id":    UserID,
		"product_id": ProductID,
	}

	// Corrected update logic
	update := bson.M{
		"$set": bson.M{
			"rating": Rating, // Always update the rating
		},
		"$setOnInsert": bson.M{
			"_id":        id, // Only set these fields when a new document is created
			"user_id":    UserID,
			"product_id": ProductID,
		},
	}

	_, err := d.collection.UpdateOne(
		ctx,
		filter,
		update,
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return richerror.New(op).WithWarpError(err)
	}

	return nil
}

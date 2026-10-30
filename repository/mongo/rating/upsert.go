package mongodbrating

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Upsert stores the user's rating and returns the product's new average
// rating and how many users have rated it.
func (d DB) Upsert(ctx context.Context, UserID string, ProductID string, Rating float64) (float64, int, error) {
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
		return 0, 0, richerror.New(op).WithWarpError(err)
	}
	avgFilter := bson.M{"product_id": ProductID}

	// --- FIX: Use Find() to get a cursor for all matching documents, not FindOne() ---
	cursor, err := d.collection.Find(ctx, avgFilter)
	if err != nil {
		return 0, 0, richerror.New(op).WithWarpError(err)
	}
	defer cursor.Close(ctx) // Ensure the cursor is closed when the function exits

	var ratings []entity.Rating
	// --- FIX: Decode all documents from the cursor into the slice ---
	if err = cursor.All(ctx, &ratings); err != nil {
		return 0, 0, richerror.New(op).WithWarpError(err)
	}

	// The rest of your averaging logic is correct
	if len(ratings) == 0 {
		return 0, 0, nil // No ratings found for this product, average is 0.
	}

	var totalRating float64
	for _, r := range ratings {
		totalRating += r.Rating
	}

	avgRating := totalRating / float64(len(ratings))
	return avgRating, len(ratings), nil
}

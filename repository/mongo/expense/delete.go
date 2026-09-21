package mongodbexpense

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// Delete soft-deletes an expense: accounting rows are never physically
// removed, so a past report stays reproducible.
func (d DB) Delete(ctx context.Context, id string) error {
	const op = "mongodbexpense.Delete"

	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now(), "updated_at": time.Now()}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("expense.not_found")
	}
	return nil
}

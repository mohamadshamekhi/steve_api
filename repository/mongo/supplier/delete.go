package mongodbsupplier

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// Delete soft-deletes a supplier: purchase invoices keep pointing at it, so
// the row must stay readable for accounting history.
func (d DB) Delete(ctx context.Context, id string) error {
	const op = "mongodbsupplier.Delete"

	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now(), "is_active": false}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("supplier.not_found")
	}
	return nil
}

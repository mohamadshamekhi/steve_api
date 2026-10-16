package mongodbrepair

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// CountOpen counts repair requests that are still in the workshop's hands:
// anything not completed, shipped back or cancelled.
func (d DB) CountOpen(ctx context.Context) (int64, error) {
	const op = "mongodbrepair.CountOpen"

	count, err := d.collection.CountDocuments(ctx, bson.M{
		"deleted_at": bson.M{"$exists": false},
		"status": bson.M{"$nin": bson.A{
			entity.StatusCompleted, entity.StatusShippedBack, entity.StatusCanceled,
		}},
	})
	if err != nil {
		return 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return count, nil
}

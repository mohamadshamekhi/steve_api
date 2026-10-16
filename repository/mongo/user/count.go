package mongodbuser

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// CountCreatedBetween counts customers who signed up inside the window.
func (d DB) CountCreatedBetween(ctx context.Context, from, to time.Time) (int64, error) {
	const op = "mongodbuser.CountCreatedBetween"

	count, err := d.collection.CountDocuments(ctx, bson.M{
		"created_at": bson.M{"$gte": from, "$lte": to},
		"deleted_at": bson.M{"$exists": false},
	})
	if err != nil {
		return 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return count, nil
}

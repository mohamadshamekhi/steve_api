package mongodbsupplier

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) Update(ctx context.Context, s entity.Supplier) error {
	const op = "mongodbsupplier.Update"

	s.UpdatedAt = time.Now()
	filter := bson.M{"_id": s.ID, "deleted_at": bson.M{"$exists": false}}

	result, err := d.collection.ReplaceOne(ctx, filter, s)
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

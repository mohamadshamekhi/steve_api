package mongodbexpense

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) Update(ctx context.Context, e entity.Expense) error {
	const op = "mongodbexpense.Update"

	e.UpdatedAt = time.Now()
	filter := bson.M{"_id": e.ID, "deleted_at": bson.M{"$exists": false}}

	result, err := d.collection.ReplaceOne(ctx, filter, e)
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

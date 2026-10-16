package mongodbsalesreturn

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) Update(ctx context.Context, r entity.SalesReturn) error {
	const op = "mongodbsalesreturn.Update"

	r.UpdatedAt = time.Now()

	result, err := d.collection.ReplaceOne(ctx, bson.M{"_id": r.ID}, r)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("sales_return.not_found")
	}
	return nil
}

package mongodbsell

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) Update(ctx context.Context, sell entity.Sell) error {
	const op = "mongodbsell.Update"
	filter := bson.M{"_id": sell.ID}
	result, err := d.collection.ReplaceOne(ctx, filter, sell)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found")
	}

	return nil
}

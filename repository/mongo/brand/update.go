package mongodbbrand

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) Update(ctx context.Context, brand entity.Brand) error {
	const op = "mongodbbrand"
	filter := bson.M{"_id": brand.ID}
	result, err := d.collection.ReplaceOne(ctx, filter, brand)
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

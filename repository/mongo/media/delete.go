package mongodbmedia

import (
	"context"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) DeleteMedia(ctx context.Context, id string) error {
	const op = "mongodb.DeleteMedia"

	filter := bson.M{"_id": id}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").
			WithKind(richerror.KindUnexpected)
	}

	if result.DeletedCount == 0 {
		return richerror.New(op).
			WithMessageKey("media.not_found").
			WithKind(richerror.KindNotFound).
			WithMeta(map[string]any{"id": id})
	}

	return nil
}

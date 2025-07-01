package mongodbmedia

import (
	"context"
	"errors"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) GetMediaByID(ctx context.Context, id string) (entity.Media, error) {
	const op = "mongodb.getMediaByID"
	filter := bson.M{"_id": id}
	var media entity.Media

	if err := d.collection.FindOne(ctx, filter).Decode(&media); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return media, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return media, richerror.New(op).WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}
	return media, nil
}

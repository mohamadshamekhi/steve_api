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
func (d DB) FindManyByIDs(ctx context.Context, ids []string) ([]entity.Media, error) {
	const op = "mongodb.FindManyByIDs"

	if len(ids) == 0 {
		return []entity.Media{}, nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var mediaList []entity.Media
	if err := cursor.All(ctx, &mediaList); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return mediaList, nil
}

package mongodbcomment

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Comment, error) {
	const op = "mongodbcomment.FindByID"
	filter := bson.M{"_id": id}
	var cm entity.Comment

	if err := d.collection.FindOne(ctx, filter).Decode(&cm); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return cm, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return cm, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}
	return cm, nil
}

func (d DB) FindManyByIDs(ctx context.Context, ids []string) (map[string]entity.Comment, error) {
	const op = "mongodbcomment.FindManyByIDs"

	if len(ids) == 0 {
		return make(map[string]entity.Comment), nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	resultMap := make(map[string]entity.Comment)

	for cursor.Next(ctx) {
		var comment entity.Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, richerror.New(op).WithWarpError(err).
				WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
		}
		resultMap[comment.ID] = comment
	}

	if err := cursor.Err(); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return resultMap, nil
}

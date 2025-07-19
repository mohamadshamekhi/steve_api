package mongodbproduct

import (
	"context"
	"errors"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Product, error) {
	const op = "mongodbproduct.FindByID"
	filter := bson.M{
		"$or": []bson.M{
			{"_id": id},
		},
		"deleted_at": bson.M{"$exists": false},
	}
	var p entity.Product

	if err := d.collection.FindOne(ctx, filter).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return p, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}

	return p, nil
}

func (d DB) FindBySlug(ctx context.Context, slug string) (entity.Product, error) {
	const op = "mongodbproduct.FindBySlug"
	filter := bson.M{
		"$or": []bson.M{
			{"slug": slug},
		},
		"deleted_at": bson.M{"$exists": false},
	}
	var p entity.Product

	if err := d.collection.FindOne(ctx, filter).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return p, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"slug": slug})
	}

	return p, nil
}

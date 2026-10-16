package mongodbbrand

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Brand, error) {
	const op = "mongodbbrand.FindByID"
	filter := bson.M{"_id": id}
	var brand entity.Brand

	if err := d.collection.FindOne(ctx, filter).Decode(&brand); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return brand, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return brand, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}
	return brand, nil
}

func (d DB) FindBySlug(ctx context.Context, slug string) (entity.Brand, error) {
	const op = "mongodbbrand.FindBySlug"
	filter := bson.M{"slug": slug}
	var brand entity.Brand

	if err := d.collection.FindOne(ctx, filter).Decode(&brand); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return brand, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return brand, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"slug": slug})
	}
	return brand, nil
}

func (d DB) FindByName(ctx context.Context, name string) (entity.Brand, error) {
	const op = "mongodbbrand.FindByName"
	filter := bson.M{"name": name}
	var brand entity.Brand

	if err := d.collection.FindOne(ctx, filter).Decode(&brand); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return brand, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return brand, richerror.New(op).WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"name": name})
	}
	return brand, nil
}

func (d DB) GetBrandsByIDs(ctx context.Context, ids []string) (map[string]entity.Brand, error) {
	const op = "mongodbbrand.GetBrandsByIDs"

	if len(ids) == 0 {
		return make(map[string]entity.Brand), nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	resultMap := make(map[string]entity.Brand)

	for cursor.Next(ctx) {
		var brand entity.Brand
		if err := cursor.Decode(&brand); err != nil {
			return nil, richerror.New(op).WithWarpError(err).
				WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
		}
		resultMap[brand.ID] = brand
	}

	if err := cursor.Err(); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return resultMap, nil
}

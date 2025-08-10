package mongodbproduct

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (d DB) FindManyByIDs(ctx context.Context, ids []string) ([]entity.Product, error) {
	const op = "mongodbproduct.FindManyByIDs"
	if len(ids) == 0 {
		return []entity.Product{}, nil
	}

	filter := bson.M{
		"_id":        bson.M{"$in": ids},
		"deleted_at": bson.M{"$exists": false},
	}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var products []entity.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return products, nil
}

func (d DB) FindNonExistingIDs(ctx context.Context, ids []string) ([]string, error) {
	const op = "mongodbproduct.FindNonExistingIDs"

	if len(ids) == 0 {
		return nil, nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	opts := options.Find().SetProjection(bson.M{"_id": 1})

	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID string `bson:"_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	existingIDSet := make(map[string]struct{}, len(results))
	for _, result := range results {
		existingIDSet[result.ID] = struct{}{}
	}

	var notFoundIDs []string
	for _, id := range ids {
		if _, found := existingIDSet[id]; !found {
			notFoundIDs = append(notFoundIDs, id)
		}
	}

	return notFoundIDs, nil

}

package mongodbcategory

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

func (d DB) FindByID(ctx context.Context, id string) (entity.Category, error) {
	const op = "mongodbcategory.FindByID"
	filter := bson.M{"_id": id}
	var category entity.Category

	if err := d.collection.FindOne(ctx, filter).Decode(&category); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return category, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return category, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}

	return category, nil
}

func (d DB) FindBySlug(ctx context.Context, slug string) (entity.Category, error) {
	const op = "mongodbcategory.FindByID"
	filter := bson.M{"slug": slug}
	var category entity.Category

	if err := d.collection.FindOne(ctx, filter).Decode(&category); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return category, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return category, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"slug": slug})
	}

	return category, nil
}

func (d DB) GetCategoriesByIDs(ctx context.Context, ids []string) (map[string]entity.Category, error) {
	const op = "mongodbcategory.GetCategoriesByIDs"
	if len(ids) == 0 {
		return make(map[string]entity.Category), nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	resultMap := make(map[string]entity.Category)

	for cursor.Next(ctx) {
		var cate entity.Category
		if err := cursor.Decode(&cate); err != nil {
			return nil, richerror.New(op).WithWarpError(err).
				WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
		}
		resultMap[cate.ID] = cate
	}

	if err := cursor.Err(); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return resultMap, nil
}

func (d DB) FindNonExistingIDs(ctx context.Context, ids []string) ([]string, error) {
	const op = "mongodbcategory.FindNonExistingIDs"

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

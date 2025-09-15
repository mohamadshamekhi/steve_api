package mongodbarticle

import (
	"context"
	"errors"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Article, error) {
	const op = "mongodbarticle.FindByID"
	filter := bson.M{"_id": id}
	var article entity.Article
	fmt.Printf("DEBUG: Executing FindOne with filter: %v\n", filter)

	if err := d.collection.FindOne(ctx, filter).Decode(&article); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return article, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return article, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}

	return article, nil
}

func (d DB) FindBySlug(ctx context.Context, slug string) (entity.Article, error) {
	const op = "mongodbarticle.FindBySlug"
	filter := bson.M{"slug": slug}
	var article entity.Article

	if err := d.collection.FindOne(ctx, filter).Decode(&article); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return article, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return article, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"slug": slug})
	}

	return article, nil
}

func (d DB) FindManyByIDs(ctx context.Context, ids []string) ([]entity.Article, error) {
	const op = "mongodbarticle.FindManyByIDs"
	if len(ids) == 0 {
		return []entity.Article{}, nil
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

	var articles []entity.Article
	if err := cursor.All(ctx, &articles); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return articles, nil
}

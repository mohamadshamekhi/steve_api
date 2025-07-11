package mongodbarticle

import (
	"context"
	"errors"
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
	const op = "mongodbarticle.FindByID"
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

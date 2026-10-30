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

// FindByPreviousSlug finds the published article that used to live at slug
// before a rename.
func (d DB) FindByPreviousSlug(ctx context.Context, slug string) (entity.Article, error) {
	const op = "mongodbarticle.FindByPreviousSlug"
	filter := bson.M{
		"previous_slugs": slug,
		"is_published":   true,
		"deleted_at":     bson.M{"$exists": false},
	}
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

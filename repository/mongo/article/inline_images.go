package mongodbarticle

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindWithInlineImages returns every non-deleted article, published or not,
// whose body still embeds a base64 image.
func (d DB) FindWithInlineImages(ctx context.Context) ([]entity.Article, error) {
	const op = "mongodbarticle.FindWithInlineImages"
	filter := bson.M{
		"body":       primitive.Regex{Pattern: `src\s*=\s*["']data:image/`, Options: "i"},
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

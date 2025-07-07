package mongodbcategory

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (d DB) Delete(ctx context.Context, category entity.Category) error {
	const op = "mongodbcategory.Delete"
	category.DeletedAt = time.Now()
	filter := bson.M{"_id": category.ID}
	result, err := d.collection.ReplaceOne(ctx, filter, category)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found")
	}

	return nil
}

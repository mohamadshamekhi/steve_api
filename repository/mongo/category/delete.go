package mongodbcategory

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) Delete(ctx context.Context, category entity.Category) error {
	const op = "mongodbcategory.Delete"
	category.DeletedAt = time.Now()
	filter := bson.M{
		"$or": []bson.M{
			{"_id": category.ID},
			{"ancestors._id": category.ID},
		},
		"deleted_at": bson.M{"$exists": false},
	}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	res, err := d.collection.UpdateMany(ctx, filter, update)

	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if res.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found")
	}

	return nil
}

package mongodbproduct

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (d DB) Update(ctx context.Context, p entity.Product) error {
	const op = "mongodbproduct.Update"

	filter := bson.M{
		"$or": []bson.M{
			{"_id": p.ID},
		},
		"deleted_at": bson.M{"$exists": false},
	}
	p.UpdatedAt = time.Now()
	result, err := d.collection.ReplaceOne(ctx, filter, p)
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

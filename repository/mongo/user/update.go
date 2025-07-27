package mongodbuser

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (d DB) UpdateUser(ctx context.Context, u entity.User) error {
	const op = "mongodb.UpdateUser"
	filter := bson.M{"_id": u.ID}
	u.UpdatedAt = time.Now()
	result, err := d.collection.ReplaceOne(ctx, filter, u)
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

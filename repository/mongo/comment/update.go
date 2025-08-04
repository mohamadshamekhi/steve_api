package mongodbcomment

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (d DB) Update(ctx context.Context, comment entity.Comment) error {
	const op = "mongodbcomment.Update"

	filter := bson.M{"_id": comment.ID}
	comment.UpdatedAt = time.Now()
	result, err := d.collection.ReplaceOne(ctx, filter, comment)
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

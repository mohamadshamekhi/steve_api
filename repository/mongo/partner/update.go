package mongodbpartner

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) Update(ctx context.Context, p entity.Partner) error {
	const op = "mongodbpartner.Update"
	p.UpdatedAt = time.Now()
	result, err := d.collection.ReplaceOne(ctx, bson.M{"_id": p.ID}, p)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("partner.not_found")
	}
	return nil
}

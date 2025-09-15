package mongodbdiscount

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) Update(ctx context.Context, discount entity.Discount) error {
	const op = "mongodbdiscount.Update"
	filter := bson.M{"_id": discount.ID}
	discount.UpdatedAt = time.Now()
	result, err := d.collection.ReplaceOne(ctx, filter, discount)
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

func (d DB) UpdateUsageCount(ctx context.Context, discountID string) error {
	const op = "mongodbdiscount.UpdateUsageCount"
	filter := bson.M{"_id": discountID}
	update := bson.M{
		"$inc": bson.M{"usage_count": 1},
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("discount.not_found").
			WithMeta(map[string]any{"discount_id": discountID})
	}

	return nil

}

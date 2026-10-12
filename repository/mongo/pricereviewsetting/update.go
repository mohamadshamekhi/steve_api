package mongodbpricereviewsetting

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) Update(ctx context.Context, intervalDays int) error {
	const op = "mongodbpricereviewsetting.Update"

	filter := bson.M{"_id": entity.PriceReviewSettingID}
	update := bson.M{
		"$set": bson.M{
			"interval_days": intervalDays,
			"updated_at":    time.Now(),
		},
	}

	_, err := d.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	return nil
}

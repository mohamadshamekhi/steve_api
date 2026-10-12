package mongodbpricereviewsetting

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) Get(ctx context.Context) (entity.PriceReviewSetting, error) {
	const op = "mongodbpricereviewsetting.Get"

	var s entity.PriceReviewSetting
	err := d.collection.FindOne(ctx, bson.M{"_id": entity.PriceReviewSettingID}).Decode(&s)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.PriceReviewSetting{
				ID:           entity.PriceReviewSettingID,
				IntervalDays: entity.DefaultPriceReviewIntervalDays,
			}, nil
		}
		return entity.PriceReviewSetting{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	return s, nil
}

package mongodborder

import (
	"context"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// ClaimThankYouCoupon atomically flips thank_you_claimed; only the caller
// that gets true may mint the coupon.
func (d DB) ClaimThankYouCoupon(ctx context.Context, orderID string) (bool, error) {
	const op = "mongodborder.ClaimThankYouCoupon"

	filter := bson.M{"_id": orderID, "thank_you_claimed": bson.M{"$ne": true}}
	res, err := d.collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"thank_you_claimed": true}})
	if err != nil {
		return false, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	return res.ModifiedCount == 1, nil
}

func (d DB) ReleaseThankYouCoupon(ctx context.Context, orderID string) error {
	const op = "mongodborder.ReleaseThankYouCoupon"

	_, err := d.collection.UpdateOne(ctx, bson.M{"_id": orderID}, bson.M{"$unset": bson.M{"thank_you_claimed": ""}})
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	return nil
}

func (d DB) SetThankYouDiscountID(ctx context.Context, orderID string, discountID string) error {
	const op = "mongodborder.SetThankYouDiscountID"

	_, err := d.collection.UpdateOne(ctx, bson.M{"_id": orderID}, bson.M{"$set": bson.M{"thank_you_discount_id": discountID}})
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	return nil
}

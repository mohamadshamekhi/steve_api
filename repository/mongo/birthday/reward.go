package mongodbbirthday

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ClaimReward inserts the reward row before any coupon is minted. ok is false
// when the user already has a reward for that year (unique index), so two
// scheduler runs can never grant twice.
func (d DB) ClaimReward(ctx context.Context, r entity.BirthdayReward) (bool, error) {
	const op = "mongodbbirthday.ClaimReward"
	if _, err := d.rewards.InsertOne(ctx, r); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return false, nil
		}
		return false, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return true, nil
}

// ReleaseReward deletes a claim whose coupon could not be minted, so the next
// run retries it.
func (d DB) ReleaseReward(ctx context.Context, id string) error {
	const op = "mongodbbirthday.ReleaseReward"
	if _, err := d.rewards.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return nil
}

func (d DB) SetRewardDiscount(ctx context.Context, id, discountID, discountCode string) error {
	const op = "mongodbbirthday.SetRewardDiscount"
	_, err := d.rewards.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"discount_id":   discountID,
		"discount_code": discountCode,
	}})
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return nil
}

func (d DB) MarkRewardSMSSent(ctx context.Context, id string) error {
	const op = "mongodbbirthday.MarkRewardSMSSent"
	if _, err := d.rewards.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"sms_sent": true}}); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return nil
}

func (d DB) ListRewards(ctx context.Context, p commonparam.PaginationRequest) ([]entity.BirthdayReward, int64, error) {
	const op = "mongodbbirthday.ListRewards"
	filter := bson.M{}

	total, err := d.rewards.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(p.GetOffset())).
		SetLimit(int64(p.GetPageSize()))
	cursor, err := d.rewards.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var out []entity.BirthdayReward
	if err := cursor.All(ctx, &out); err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return out, total, nil
}

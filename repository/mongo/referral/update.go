package mongodbreferral

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// MarkReferralJoined promotes an outstanding SMS invite to "joined".
func (d DB) MarkReferralJoined(ctx context.Context, referralID, referredUserID string) error {
	const op = "mongodbreferral.MarkReferralJoined"
	now := time.Now()
	update := bson.M{"$set": bson.M{
		"status":           entity.ReferralStatusJoined,
		"referred_user_id": referredUserID,
		"joined_at":        now,
		"updated_at":       now,
	}}
	if _, err := d.referrals.UpdateOne(ctx, bson.M{"_id": referralID}, update); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return nil
}

// MarkReferralsRewarded stamps the given referrals with the reward that consumed
// them so they are not counted again, but only those not already stamped. It
// returns how many rows it actually claimed, so a caller racing another reward
// grant can detect that it lost and back out.
func (d DB) MarkReferralsRewarded(ctx context.Context, referralIDs []string, rewardBatchID string) (int64, error) {
	const op = "mongodbreferral.MarkReferralsRewarded"
	if len(referralIDs) == 0 {
		return 0, nil
	}
	filter := bson.M{
		"_id": bson.M{"$in": referralIDs},
		"$or": []bson.M{
			{"reward_batch_id": bson.M{"$exists": false}},
			{"reward_batch_id": ""},
		},
	}
	update := bson.M{"$set": bson.M{"reward_batch_id": rewardBatchID, "updated_at": time.Now()}}
	res, err := d.referrals.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return res.ModifiedCount, nil
}

// UnmarkReferralsRewarded reverses MarkReferralsRewarded for a batch id — used to
// roll back when the reward could not be finalised.
func (d DB) UnmarkReferralsRewarded(ctx context.Context, rewardBatchID string) error {
	const op = "mongodbreferral.UnmarkReferralsRewarded"
	update := bson.M{"$set": bson.M{"reward_batch_id": "", "updated_at": time.Now()}}
	if _, err := d.referrals.UpdateMany(ctx, bson.M{"reward_batch_id": rewardBatchID}, update); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return nil
}

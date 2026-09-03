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

// MarkReferralsRewarded stamps a batch of referrals with the reward that consumed
// them so they are not counted again.
func (d DB) MarkReferralsRewarded(ctx context.Context, referralIDs []string, rewardBatchID string) error {
	const op = "mongodbreferral.MarkReferralsRewarded"
	if len(referralIDs) == 0 {
		return nil
	}
	update := bson.M{"$set": bson.M{"reward_batch_id": rewardBatchID, "updated_at": time.Now()}}
	if _, err := d.referrals.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": referralIDs}}, update); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return nil
}

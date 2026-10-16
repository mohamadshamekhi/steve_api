package mongodbreferral

import (
	"context"
	"errors"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindAccountByUserID(ctx context.Context, userID string) (entity.ReferralAccount, error) {
	const op = "mongodbreferral.FindAccountByUserID"
	var a entity.ReferralAccount
	if err := d.accounts.FindOne(ctx, bson.M{"user_id": userID}).Decode(&a); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return a, richerror.New(op).WithWarpError(err).
				WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return a, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return a, nil
}

func (d DB) FindAccountByCode(ctx context.Context, code string) (entity.ReferralAccount, error) {
	const op = "mongodbreferral.FindAccountByCode"
	var a entity.ReferralAccount
	if err := d.accounts.FindOne(ctx, bson.M{"code": code}).Decode(&a); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return a, richerror.New(op).WithWarpError(err).
				WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return a, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return a, nil
}

func (d DB) CreateAccount(ctx context.Context, a entity.ReferralAccount) error {
	const op = "mongodbreferral.CreateAccount"
	if _, err := d.accounts.InsertOne(ctx, a); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return richerror.New(op).WithWarpError(err).
				WithKind(richerror.KindBadRequest).WithMessageKey("referral.account_conflict")
		}
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return nil
}

// TryConsumeRewardBatch atomically reserves one reward batch: it bumps
// rewarded_count by threshold only if the account currently has at least
// `threshold` joined-but-unrewarded invitees, and returns the updated account.
// ok is false when there is no full batch to consume — e.g. a concurrent signup
// hook already reserved it — so the caller must not grant a coupon.
func (d DB) TryConsumeRewardBatch(ctx context.Context, userID string, threshold int) (entity.ReferralAccount, bool, error) {
	const op = "mongodbreferral.TryConsumeRewardBatch"
	filter := bson.M{
		"user_id": userID,
		"$expr": bson.M{"$gte": bson.A{
			bson.M{"$subtract": bson.A{"$total_joined", "$rewarded_count"}}, threshold,
		}},
	}
	update := bson.M{
		"$inc": bson.M{"rewarded_count": threshold},
		"$set": bson.M{"updated_at": time.Now()},
	}
	var a entity.ReferralAccount
	err := d.accounts.FindOneAndUpdate(ctx, filter, update, mongoFindOneAndUpdateAfter()).Decode(&a)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return a, false, nil
	}
	if err != nil {
		return a, false, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return a, true, nil
}

// ReleaseRewardBatch rolls back a TryConsumeRewardBatch reservation when the
// caller could not finalise the reward.
func (d DB) ReleaseRewardBatch(ctx context.Context, userID string, threshold int) error {
	_, err := d.IncrementAccountCounters(ctx, userID, 0, -threshold)
	return err
}

// IncrementAccountCounters atomically bumps total_joined and rewarded_count by
// the given deltas and returns the updated document.
func (d DB) IncrementAccountCounters(ctx context.Context, userID string, joinedDelta, rewardedDelta int) (entity.ReferralAccount, error) {
	const op = "mongodbreferral.IncrementAccountCounters"
	update := bson.M{
		"$inc": bson.M{"total_joined": joinedDelta, "rewarded_count": rewardedDelta},
		"$set": bson.M{"updated_at": time.Now()},
	}
	opts := mongoFindOneAndUpdateAfter()
	var a entity.ReferralAccount
	if err := d.accounts.FindOneAndUpdate(ctx, bson.M{"user_id": userID}, update, opts).Decode(&a); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return a, richerror.New(op).WithWarpError(err).
				WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return a, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return a, nil
}

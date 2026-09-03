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

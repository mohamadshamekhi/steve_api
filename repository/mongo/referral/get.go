package mongodbreferral

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func mongoFindOneAndUpdateAfter() *options.FindOneAndUpdateOptions {
	return options.FindOneAndUpdate().SetReturnDocument(options.After)
}

func (d DB) FindReferralByReferredUserID(ctx context.Context, referredUserID string) (entity.Referral, error) {
	const op = "mongodbreferral.FindReferralByReferredUserID"
	var r entity.Referral
	if err := d.referrals.FindOne(ctx, bson.M{"referred_user_id": referredUserID}).Decode(&r); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return r, richerror.New(op).WithWarpError(err).
				WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return r, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return r, nil
}

// FindPendingInviteByPhone returns the newest outstanding SMS invite for a phone.
func (d DB) FindPendingInviteByPhone(ctx context.Context, phone string) (entity.Referral, error) {
	const op = "mongodbreferral.FindPendingInviteByPhone"
	var r entity.Referral
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	filter := bson.M{"invited_phone": phone, "status": entity.ReferralStatusInvited}
	if err := d.referrals.FindOne(ctx, filter, opts).Decode(&r); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return r, richerror.New(op).WithWarpError(err).
				WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return r, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return r, nil
}

// ExistsInviteByPhone reports whether any referral row (from any referrer, in any
// status) has already targeted this phone number.
func (d DB) ExistsInviteByPhone(ctx context.Context, phone string) (bool, error) {
	const op = "mongodbreferral.ExistsInviteByPhone"
	err := d.referrals.FindOne(ctx, bson.M{"invited_phone": phone}).Err()
	if err == nil {
		return true, nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	return false, richerror.New(op).WithWarpError(err).
		WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
}

// ListUnrewardedJoined returns joined referrals with no reward batch yet, oldest
// first, capped at limit.
func (d DB) ListUnrewardedJoined(ctx context.Context, referrerUserID string, limit int) ([]entity.Referral, error) {
	const op = "mongodbreferral.ListUnrewardedJoined"
	filter := bson.M{
		"referrer_user_id": referrerUserID,
		"status":           entity.ReferralStatusJoined,
		"$or": []bson.M{
			{"reward_batch_id": bson.M{"$exists": false}},
			{"reward_batch_id": ""},
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "joined_at", Value: 1}}).SetLimit(int64(limit))
	cursor, err := d.referrals.Find(ctx, filter, opts)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)
	var out []entity.Referral
	if err := cursor.All(ctx, &out); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return out, nil
}

func (d DB) ListReferralsByReferrer(ctx context.Context, referrerUserID string, p commonparam.PaginationRequest) ([]entity.Referral, int64, error) {
	const op = "mongodbreferral.ListReferralsByReferrer"
	filter := bson.M{"referrer_user_id": referrerUserID}

	total, err := d.referrals.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(p.GetOffset())).
		SetLimit(int64(p.GetPageSize()))
	cursor, err := d.referrals.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)
	var out []entity.Referral
	if err := cursor.All(ctx, &out); err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return out, total, nil
}

func (d DB) ListRewardsByReferrer(ctx context.Context, referrerUserID string) ([]entity.ReferralReward, error) {
	const op = "mongodbreferral.ListRewardsByReferrer"
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := d.rewards.Find(ctx, bson.M{"referrer_user_id": referrerUserID}, opts)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)
	var out []entity.ReferralReward
	if err := cursor.All(ctx, &out); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return out, nil
}

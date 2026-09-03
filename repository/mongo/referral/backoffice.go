package mongodbreferral

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AllAccounts returns referral accounts sorted by total_joined desc (the
// leaderboard), optionally filtered to a minimum join count.
func (d DB) AllAccounts(ctx context.Context, minJoined int, p commonparam.PaginationRequest) ([]entity.ReferralAccount, int64, error) {
	const op = "mongodbreferral.AllAccounts"

	filter := bson.M{}
	if minJoined > 0 {
		filter["total_joined"] = bson.M{"$gte": minJoined}
	}

	total, err := d.accounts.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "total_joined", Value: -1}, {Key: "created_at", Value: -1}}).
		SetSkip(int64(p.GetOffset())).
		SetLimit(int64(p.GetPageSize()))

	cursor, err := d.accounts.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var out []entity.ReferralAccount
	if err := cursor.All(ctx, &out); err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return out, total, nil
}

// CountRewardsByReferrer is used to fill RewardsEarned on a leaderboard row.
func (d DB) CountRewardsByReferrer(ctx context.Context, referrerUserID string) (int64, error) {
	const op = "mongodbreferral.CountRewardsByReferrer"
	n, err := d.rewards.CountDocuments(ctx, bson.M{"referrer_user_id": referrerUserID})
	if err != nil {
		return 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return n, nil
}

// Stats returns program-wide counters for the admin dashboard.
func (d DB) Stats(ctx context.Context) (accounts, invites, joined, rewards int64, err error) {
	const op = "mongodbreferral.Stats"

	if accounts, err = d.accounts.CountDocuments(ctx, bson.M{}); err != nil {
		return 0, 0, 0, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if invites, err = d.referrals.CountDocuments(ctx, bson.M{}); err != nil {
		return 0, 0, 0, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if joined, err = d.referrals.CountDocuments(ctx, bson.M{"status": entity.ReferralStatusJoined}); err != nil {
		return 0, 0, 0, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if rewards, err = d.rewards.CountDocuments(ctx, bson.M{}); err != nil {
		return 0, 0, 0, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return accounts, invites, joined, rewards, nil
}

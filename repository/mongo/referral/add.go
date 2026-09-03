package mongodbreferral

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) CreateReferral(ctx context.Context, r entity.Referral) error {
	const op = "mongodbreferral.CreateReferral"
	if _, err := d.referrals.InsertOne(ctx, r); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"referral": r})
	}
	return nil
}

func (d DB) CreateReward(ctx context.Context, rw entity.ReferralReward) error {
	const op = "mongodbreferral.CreateReward"
	if _, err := d.rewards.InsertOne(ctx, rw); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"reward": rw})
	}
	return nil
}

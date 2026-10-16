package mongodbreferral

import (
	"context"
	"errors"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetSettings returns the single referral-settings document, or KindNotFound
// when an admin has never customised it (the service then uses config defaults).
func (d DB) GetSettings(ctx context.Context) (entity.ReferralSettings, error) {
	const op = "mongodbreferral.GetSettings"
	var s entity.ReferralSettings
	err := d.settings.FindOne(ctx, bson.M{"_id": entity.ReferralSettingsID}).Decode(&s)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return s, richerror.New(op).WithWarpError(err).
				WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return s, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return s, nil
}

// UpsertSettings replaces the single settings document.
func (d DB) UpsertSettings(ctx context.Context, s entity.ReferralSettings) error {
	const op = "mongodbreferral.UpsertSettings"
	s.ID = entity.ReferralSettingsID
	s.UpdatedAt = time.Now()
	_, err := d.settings.ReplaceOne(ctx, bson.M{"_id": entity.ReferralSettingsID}, s,
		options.Replace().SetUpsert(true))
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return nil
}

package mongodbreferral

import (
	"context"

	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	accountsCollectionName  = "referral_accounts"
	referralsCollectionName = "referrals"
	rewardsCollectionName   = "referral_rewards"
	settingsCollectionName  = "referral_settings"
)

// DB holds the collections the referral domain persists to.
type DB struct {
	accounts  *mongo.Collection
	referrals *mongo.Collection
	rewards   *mongo.Collection
	settings  *mongo.Collection
}

func New(c *mongodb.DB) *DB {
	accounts := c.Connection().Collection(accountsCollectionName)
	referrals := c.Connection().Collection(referralsCollectionName)
	rewards := c.Connection().Collection(rewardsCollectionName)
	settings := c.Connection().Collection(settingsCollectionName)

	if _, err := accounts.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "code", Value: 1}}, Options: options.Index().SetUnique(true)},
	}); err != nil {
		panic(err)
	}

	// Old rows (written before "one invite per phone" existed) can hold duplicate
	// invited_phone values; collapse them before the unique index is built or the
	// build panics.
	if err := dedupeInvitedPhones(context.Background(), referrals); err != nil {
		panic(err)
	}

	if _, err := referrals.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "referrer_user_id", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "referred_user_id", Value: 1}}},
		{Keys: bson.D{{Key: "invited_phone", Value: 1}, {Key: "status", Value: 1}}},
		// One invite per phone number, ever (link-source rows have no
		// invited_phone and are excluded by the partial filter).
		{
			Keys: bson.D{{Key: "invited_phone", Value: 1}},
			Options: options.Index().SetUnique(true).
				SetPartialFilterExpression(bson.M{"invited_phone": bson.M{"$type": "string"}}),
		},
	}); err != nil {
		panic(err)
	}

	if _, err := rewards.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "referrer_user_id", Value: 1}, {Key: "created_at", Value: -1}},
	}); err != nil {
		panic(err)
	}

	return &DB{accounts: accounts, referrals: referrals, rewards: rewards, settings: settings}
}

// dedupeInvitedPhones keeps a single referral row per non-empty invited_phone —
// preferring a "joined" row, then the most recent — and deletes the rest, so the
// unique partial index on invited_phone can be built on legacy data.
func dedupeInvitedPhones(ctx context.Context, referrals *mongo.Collection) error {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"invited_phone": bson.M{"$type": "string", "$ne": ""}}}},
		// joined before invited, then newest first — $first is the row to keep.
		{{Key: "$sort", Value: bson.D{{Key: "status", Value: -1}, {Key: "created_at", Value: -1}}}},
		{{Key: "$group", Value: bson.M{
			"_id":  "$invited_phone",
			"keep": bson.M{"$first": "$_id"},
			"ids":  bson.M{"$push": "$_id"},
		}}},
		{{Key: "$match", Value: bson.M{"$expr": bson.M{"$gt": bson.A{bson.M{"$size": "$ids"}, 1}}}}},
	}

	cursor, err := referrals.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var stale []interface{}
	for cursor.Next(ctx) {
		var row struct {
			Keep string   `bson:"keep"`
			IDs  []string `bson:"ids"`
		}
		if err := cursor.Decode(&row); err != nil {
			return err
		}
		for _, id := range row.IDs {
			if id != row.Keep {
				stale = append(stale, id)
			}
		}
	}
	if err := cursor.Err(); err != nil {
		return err
	}
	if len(stale) == 0 {
		return nil
	}

	_, err = referrals.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": stale}})
	return err
}

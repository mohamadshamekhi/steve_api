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

	if _, err := referrals.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "referrer_user_id", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "referred_user_id", Value: 1}}},
		{Keys: bson.D{{Key: "invited_phone", Value: 1}, {Key: "status", Value: 1}}},
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

package mongodbbirthday

import (
	"context"

	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	rewardsCollectionName  = "birthday_rewards"
	settingsCollectionName = "birthday_settings"
)

// DB holds the collections the birthday-discount domain persists to.
type DB struct {
	rewards  *mongo.Collection
	settings *mongo.Collection
}

func New(c *mongodb.DB) *DB {
	rewards := c.Connection().Collection(rewardsCollectionName)
	settings := c.Connection().Collection(settingsCollectionName)

	if _, err := rewards.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		// One birthday code per user per year — the claim in ClaimReward relies on it.
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "year", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	}); err != nil {
		panic(err)
	}

	return &DB{rewards: rewards, settings: settings}
}

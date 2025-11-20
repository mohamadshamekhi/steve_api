package mongodbdiscount

import (
	"context"
	"time"

	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "discounts"

type DB struct {
	collection *mongo.Collection
}

func (d DB) ExistSpecialDiscount(ctx context.Context, startsAt, expiresAt time.Time) (bool, error) {
	filter := bson.M{
		"is_special_discount": true,
		"is_active":           true,

		// چک هم‌پوشانی بازه‌ها:
		// existing.StartsAt < expiresAt
		// existing.ExpiresAt > startsAt
		"$and": []bson.M{
			{"starts_at": bson.M{"$lt": expiresAt}},
			{"expires_at": bson.M{"$gt": startsAt}},
		},
	}

	count, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func New(c *mongodb.DB) *DB {
	collection := c.Connection().Collection(collectionName)

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "code", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.M{
				"code": bson.M{
					"$type": "string",
				},
			}),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic(err)
	}
	return &DB{
		collection: collection,
	}
}

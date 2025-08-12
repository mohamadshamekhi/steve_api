package mongodbdiscount

import (
	"context"

	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "discounts"

type DB struct {
	collection *mongo.Collection
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

package mongodbproduct

import (
	"context"

	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "products"

type DB struct {
	collection *mongo.Collection
}

func New(c *mongodb.DB) *DB {
	collection := c.Connection().Collection(collectionName)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "en_title", Value: 1}},
			Options: options.Index(),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index(),
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		panic(err)
	}
	return &DB{
		collection: collection,
	}
}

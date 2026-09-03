package mongodbpartner

import (
	"context"

	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "partners"

type DB struct {
	collection *mongo.Collection
}

func New(c *mongodb.DB) *DB {
	collection := c.Connection().Collection(collectionName)

	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic(err)
	}

	return &DB{collection: collection}
}

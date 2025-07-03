package mongodbmedia

import (
	"context"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "media"

type DB struct {
	collection *mongo.Collection
}

func New(c *mongodb.DB) *DB {
	collection := c.Connection().Collection(collectionName)

	indexes := mongo.IndexModel{
		Keys: bson.D{{Key: "uploader_id", Value: 1}},
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexes)
	if err != nil {
		panic(err)
	}
	return &DB{
		collection: collection,
	}
}

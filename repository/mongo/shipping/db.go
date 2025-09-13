package mongodbshipping

import (
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "shipping"

type DB struct {
	collection *mongo.Collection
}

func New(c *mongodb.DB) *DB {
	collection := c.Connection().Collection(collectionName)

	return &DB{
		collection: collection,
	}
}

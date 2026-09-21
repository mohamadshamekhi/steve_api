package mongodbexpense

import (
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "expenses"

type DB struct {
	collection *mongo.Collection
}

func New(c *mongodb.DB) *DB {
	return &DB{collection: c.Connection().Collection(collectionName)}
}

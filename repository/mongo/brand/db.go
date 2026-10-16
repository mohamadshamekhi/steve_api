package mongodbbrand

import (
	"context"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/slug"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "brands"

type DB struct {
	collection *mongo.Collection
}

func New(c *mongodb.DB) *DB {
	collection := c.Connection().Collection(collectionName)

	// backfill slug for brands created before the field existed, so the
	// unique index below doesn't fail on startup (Mongo allows at most one
	// document with a missing/null value per unique index).
	backfillMissingSlugs(collection)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
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

// backfillMissingSlugs sets a slug for any brand document that doesn't have
// one yet. It's idempotent - once every document has a slug it's a cheap
// no-op find, so it's safe to run on every startup.
func backfillMissingSlugs(collection *mongo.Collection) {
	ctx := context.Background()
	filter := bson.M{"$or": []bson.M{
		{"slug": bson.M{"$exists": false}},
		{"slug": ""},
	}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	var brands []entity.Brand
	if err := cursor.All(ctx, &brands); err != nil {
		panic(err)
	}

	for _, b := range brands {
		base := slug.GenerateSlug(b.Name)
		if base == "" {
			base = b.ID
		}

		candidate := base
		counter := 1
		for {
			count, cErr := collection.CountDocuments(ctx, bson.M{"slug": candidate, "_id": bson.M{"$ne": b.ID}})
			if cErr != nil {
				panic(cErr)
			}
			if count == 0 {
				break
			}
			candidate = fmt.Sprintf("%s-%d", base, counter)
			counter++
		}

		if _, err := collection.UpdateOne(ctx, bson.M{"_id": b.ID}, bson.M{"$set": bson.M{"slug": candidate}}); err != nil {
			panic(err)
		}
	}
}

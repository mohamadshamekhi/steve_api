package mongodbproduct

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindWithTorobURL returns the non-deleted products that have a Torob link,
// with just the fields the competitor-price agent needs.
func (d DB) FindWithTorobURL(ctx context.Context) ([]entity.Product, error) {
	const op = "mongodbproduct.FindWithTorobURL"

	filter := bson.M{
		"torob_url":  bson.M{"$exists": true, "$ne": ""},
		"deleted_at": bson.M{"$exists": false},
	}
	findOptions := options.Find().
		SetSort(bson.D{{Key: "name", Value: 1}}).
		SetProjection(bson.M{
			"name": 1, "slug": 1, "torob_url": 1, "variants": 1,
			"sales_channels": 1, "is_published": 1, "kind": 1,
		})

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var products []entity.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	return products, nil
}

// SetCompetitorPrice writes only the four competitor price fields, so the
// price agent can never touch anything else on the product.
func (d DB) SetCompetitorPrice(ctx context.Context, productID string, minPrice, maxPrice uint64, offersCount uint, checkedAt time.Time) error {
	const op = "mongodbproduct.SetCompetitorPrice"

	filter := bson.M{"_id": productID, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{
		"competitor_min_price":    minPrice,
		"competitor_max_price":    maxPrice,
		"competitor_offers_count": offersCount,
		"competitor_checked_at":   checkedAt,
	}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("product.not_found").
			WithMeta(map[string]any{"product_id": productID})
	}
	return nil
}

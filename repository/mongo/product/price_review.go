package mongodbproduct

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindNeedingPriceReview returns non-deleted products whose price was last
// reviewed at or before cutoff (or never), most-overdue first.
func (d DB) FindNeedingPriceReview(ctx context.Context, cutoff time.Time, offset, limit uint) ([]entity.Product, int64, error) {
	const op = "mongodbproduct.FindNeedingPriceReview"

	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
		"$or": []bson.M{
			{"price_reviewed_at": bson.M{"$exists": false}},
			{"price_reviewed_at": bson.M{"$lte": cutoff}},
		},
	}

	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "price_reviewed_at", Value: 1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var products []entity.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return products, total, nil
}

// MarkPriceReviewed resets a product's price-review timer to at, without
// touching the rest of the document.
func (d DB) MarkPriceReviewed(ctx context.Context, productID string, at time.Time) error {
	const op = "mongodbproduct.MarkPriceReviewed"

	filter := bson.M{"_id": productID}
	update := bson.M{"$set": bson.M{"price_reviewed_at": at}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found").
			WithMeta(map[string]any{"product_id": productID})
	}
	return nil
}

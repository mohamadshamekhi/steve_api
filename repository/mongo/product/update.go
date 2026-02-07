package mongodbproduct

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) Update(ctx context.Context, p entity.Product) error {
	const op = "mongodbproduct.Update"

	filter := bson.M{
		"$or": []bson.M{
			{"_id": p.ID},
		},
		"deleted_at": bson.M{"$exists": false},
	}
	p.UpdatedAt = time.Now()
	result, err := d.collection.ReplaceOne(ctx, filter, p)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found")
	}

	return nil
}

func (d DB) UpdateQuantity(ctx context.Context, productID string, variant entity.Variant, finalQuantity uint) error {
	const op = "mongodbproduct.UpdateQuantity"

	filter := bson.M{"_id": productID}

	update := bson.M{"$set": bson.M{"variants.$[elem].quantity": finalQuantity}}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem.id": variant.ID}},
	})

	result, err := d.collection.UpdateOne(ctx, filter, update, arrayFilters)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("product.not_found").
			WithMeta(map[string]any{"product_id": productID})
	}

	if result.ModifiedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("product.variant_not_found").
			WithMeta(map[string]any{"product_id": productID, "variant_id": variant.ID})
	}

	return nil
}

func (d DB) UpdateProductAvrRating(ctx context.Context, ProductID string, AverageRating float64) error {
	const op = "mongodbproduct.UpdateProductAvrRating"
	filter := bson.M{
		"_id": ProductID,
	}

	update := bson.M{
		"$set": bson.M{
			"rating": AverageRating,
		},
	}

	_, err := d.collection.UpdateOne(
		ctx,
		filter,
		update,
		options.Update().SetUpsert(false),
	)

	if err != nil {
		return richerror.New(op).WithWarpError(err)
	}

	return nil
}

func (d DB) UpdateProductSalesCount(ctx context.Context, ProductID string, quantity uint) error {
	const op = "product.repo.UpdateProductSalesCount"

	filter := bson.M{"_id": ProductID}

	update := bson.M{
		"$inc": bson.M{"sales_count": quantity},
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("product.not_found").
			WithMeta(map[string]any{"product_id": ProductID})
	}

	return nil
}

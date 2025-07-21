package mongodbproduct

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) UpdateSpecifications(ctx context.Context, productID string, specs []entity.Specification) error {
	const op = "mongodbproduct.UpdateSpecifications"

	filter := bson.M{"_id": productID}
	update := bson.M{"$set": bson.M{"specifications": specs}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("product.not_found")
	}

	return nil
}

func (d DB) UpdateSingleSpecification(ctx context.Context, productID string, specKey string, newValue string) error {
	const op = "mongodbproduct.UpdateSingleSpecification"
	filter := bson.M{"_id": productID}

	update := bson.M{
		"$set": bson.M{"specifications.$[elem].value": newValue},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem.key": specKey}},
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
			WithMessageKey("product.specification_key_not_found").
			WithMeta(map[string]any{"spec_key": specKey})
	}

	return nil
}

package mongodbproduct

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) AddVariant(ctx context.Context, productID string, variant entity.Variant) error {
	const op = "mongodbproduct.AddVariant"

	variant.ID = ulid.Make().String()

	filter := bson.M{"_id": productID}
	update := bson.M{"$push": bson.M{"variants": variant}}

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

func (d DB) DeleteVariant(ctx context.Context, productID string, variantID string) error {
	const op = "mongodbproduct.DeleteVariant"

	filter := bson.M{"_id": productID}
	update := bson.M{"$pull": bson.M{"variants": bson.M{"id": variantID}}}

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

	if result.ModifiedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("product.variant_not_found").
			WithMeta(map[string]any{"product_id": productID, "variant_id": variantID})
	}

	return nil
}

func (d DB) UpdateVariant(ctx context.Context, productID string, variant entity.Variant) error {
	const op = "mongodbproduct.UpdateVariant"
	filter := bson.M{
		"_id":         productID,
		"variants.id": variant.ID,
	}

	update := bson.M{
		"$set": bson.M{"variants.$": variant},
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found").
			WithMeta(map[string]any{"product_id": productID, "variant_id": variant.ID})
	}

	return nil

}

package mongodbproduct

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) AddVariant(ctx context.Context, productID string, variant entity.Variant) error {
	const op = "mongodbproduct.AddVariant"

	var product entity.Product
	filter := bson.M{"_id": productID}
	err := d.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("product.not_found")
		}
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if product.Variants == nil {
		product.Variants = make([]entity.Variant, 0)
	}

	variant.BeforeCreate()

	product.Variants = append(product.Variants, variant)

	_, err = d.collection.ReplaceOne(ctx, filter, product)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
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

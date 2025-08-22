package mongodbproduct

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
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
	updateFilter := bson.M{"_id": productID, "specifications.key": specKey}
	updateQuery := bson.M{"$set": bson.M{"specifications.$.value": newValue}}

	result, err := d.collection.UpdateOne(ctx, updateFilter, updateQuery)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.ModifiedCount > 0 {
		return nil
	}

	addFilter := bson.M{
		"_id":                productID,
		"specifications.key": bson.M{"$ne": specKey},
	}

	newSpecification := entity.Specification{
		Key:   specKey,
		Value: newValue,
	}
	addQuery := bson.M{"$push": bson.M{"specifications": newSpecification}}

	result, err = d.collection.UpdateOne(ctx, addFilter, addQuery)
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

package mongodbproduct

import (
	"context"
	"errors"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Product, error) {
	const op = "mongodbproduct.FindByID"
	filter := bson.M{
		"$or": []bson.M{
			{"_id": id},
		},
		"deleted_at": bson.M{"$exists": false},
	}
	var p entity.Product

	if err := d.collection.FindOne(ctx, filter).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return p, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}

	return p, nil
}

func (d DB) FindBySlug(ctx context.Context, slug string) (entity.Product, error) {
	const op = "mongodbproduct.FindBySlug"
	filter := bson.M{
		"$or": []bson.M{
			{"slug": slug},
		},
		"deleted_at": bson.M{"$exists": false},
	}
	var p entity.Product

	if err := d.collection.FindOne(ctx, filter).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return p, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"slug": slug})
	}

	return p, nil
}

func (d DB) FindManyByIDs(ctx context.Context, ids []string) ([]entity.Product, error) {
	const op = "mongodbproduct.FindManyByIDs"
	if len(ids) == 0 {
		return []entity.Product{}, nil
	}

	filter := bson.M{
		"_id":        bson.M{"$in": ids},
		"deleted_at": bson.M{"$exists": false},
	}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var products []entity.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return products, nil
}

func (d DB) FindNonExistingIDs(ctx context.Context, ids []string) ([]string, error) {
	const op = "mongodbproduct.FindNonExistingIDs"

	if len(ids) == 0 {
		return nil, nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	opts := options.Find().SetProjection(bson.M{"_id": 1})

	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID string `bson:"_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	existingIDSet := make(map[string]struct{}, len(results))
	for _, result := range results {
		existingIDSet[result.ID] = struct{}{}
	}

	var notFoundIDs []string
	for _, id := range ids {
		if _, found := existingIDSet[id]; !found {
			notFoundIDs = append(notFoundIDs, id)
		}
	}

	return notFoundIDs, nil

}

func (d DB) GetVariantsByIDs(ctx context.Context, variantIDs []string) (map[string]entity.Variant, error) {
	const op = "mongodbproduct.GetVariantsByIDs"

	if len(variantIDs) == 0 {
		return make(map[string]entity.Variant), nil
	}

	filter := bson.M{"variants.id": bson.M{"$in": variantIDs}}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	resultMap := make(map[string]entity.Variant)

	for cursor.Next(ctx) {
		var product entity.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
		}

		for _, variant := range product.Variants {
			for _, requestedID := range variantIDs {
				if variant.ID == requestedID {
					resultMap[variant.ID] = variant
					break
				}
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return resultMap, nil
}

func (d DB) GetProductCategoriesForCartItems(ctx context.Context, items []entity.CartItem) (map[string][]string, error) {
	const op = "mongodbproduct.GetProductCategoriesForCartItems"

	if len(items) == 0 {
		return make(map[string][]string), nil
	}

	productIDs := make([]string, 0, len(items))
	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
	}

	filter := bson.M{"_id": bson.M{"$in": productIDs}}

	projection := bson.M{"_id": 1, "category_id": 1}
	opts := options.Find().SetProjection(projection)

	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	resultMap := make(map[string][]string)
	for cursor.Next(ctx) {
		var product struct {
			ID         string `bson:"_id"`
			CategoryID string `bson:"category_id"`
		}
		if err := cursor.Decode(&product); err != nil {
			return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
		}
		resultMap[product.ID] = []string{product.CategoryID}
	}

	if err := cursor.Err(); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return resultMap, nil
}

func (d DB) GetProductsByIDs(ctx context.Context, ids []string) ([]entity.Product, error) {
	const op = "mongodbproduct.GetProductsByIDs"
	if len(ids) == 0 {
		return []entity.Product{}, nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err)
	}
	defer cursor.Close(ctx)

	var products []entity.Product

	for cursor.Next(ctx) {
		var product entity.Product

		if err := cursor.Decode(&product); err != nil {
			return nil, richerror.New(op).WithWarpError(err)
		}

		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, richerror.New(op).WithWarpError(err)
	}

	return products, nil
}

func (d DB) FindProductsIDsByCategory(ctx context.Context, categoryID string) (string, error) {
	const op = "mongodbproduct.FindProductsIDsByCategory"
	filter := bson.M{
		"$or": []bson.M{
			{"category_id": categoryID},
		},
		"deleted_at": bson.M{"$exists": false},
	}

	var p entity.Product
	if err := d.collection.FindOne(ctx, filter).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p.ID, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		fmt.Println("err ", err)
		return p.ID, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"category_id": categoryID})
	}

	return p.ID, nil
}

func (d DB) GetFilterableAttributes(ctx context.Context, categoryID string) ([]productparam.AggregatedFilterOption, error) {
	const op = "product.repo.GetFilterableAttributes"

	pipeline := mongo.Pipeline{

		bson.D{
			{"$match", bson.M{
				"$or": []bson.M{
					{"category_id": categoryID},
					{"ancestors._id": categoryID},
				},
				"is_published": true,
				"deleted_at":   bson.M{"$exists": false},
			}},
		},
		bson.D{
			{"$unwind", "$specifications"},
		},

		bson.D{
			{"$group", bson.M{
				"_id":    "$specifications.key",
				"values": bson.M{"$addToSet": "$specifications.value"},
			}},
		},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var results []productparam.AggregatedFilterOption
	if err := cursor.All(ctx, &results); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return results, nil
}

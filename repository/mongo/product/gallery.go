package mongodbproduct

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) AddProductGallery(ctx context.Context, g entity.ProductGallery) error {
	const op = "mongodb.AddProductGallery"
	_, err := d.collection.InsertOne(ctx, g)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithKind(richerror.KindUnexpected).WithMeta(map[string]any{"gallery": g}).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (d DB) DeleteProductGallery(ctx context.Context, productID string, mediaID string) error {
	const op = "mongodb.DeleteProductGallery"
	filter := bson.M{"_id": productID}
	update := bson.M{
		"$pull": bson.M{
			"gallery": bson.M{"media_id": mediaID},
		},
	}
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

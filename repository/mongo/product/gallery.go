package mongodbproduct

import (
	"context"
	"errors"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) AddProductGallery(ctx context.Context, productID string, galleryItem entity.ProductGallery) error {
	const op = "mongodb.AddProductGallery"
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

	if product.Gallery == nil {
		product.Gallery = make([]entity.ProductGallery, 0)
	}

	if galleryItem.IsPrimary {
		for i := range product.Gallery {
			product.Gallery[i].IsPrimary = false
		}
	}

	product.Gallery = append(product.Gallery, galleryItem)

	update := bson.M{"$set": bson.M{"gallery": product.Gallery}}
	_, err = d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (d DB) DeleteProductGallery(ctx context.Context, productID string, mediaID string) error {
	const op = "mongodb.DeleteProductGallery"

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

	var wasPrimary bool
	removeIndex := -1
	for i, item := range product.Gallery {
		if item.MediaID == mediaID {
			wasPrimary = item.IsPrimary
			removeIndex = i
			break
		}
	}

	if removeIndex == -1 {
		return nil
	}

	newGallery := append(product.Gallery[:removeIndex], product.Gallery[removeIndex+1:]...)

	if wasPrimary && len(newGallery) > 0 {
		newGallery[0].IsPrimary = true
	}

	update := bson.M{"$set": bson.M{"gallery": newGallery}}
	_, err = d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return nil

}

func (d DB) SetPrimary(ctx context.Context, productID string, mediaID string) error {
	const op = "productservice.SetPrimary"

	initialCheckFilter := bson.M{
		"_id":              productID,
		"gallery.media_id": mediaID,
	}
	err := d.collection.FindOne(ctx, initialCheckFilter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("product.media_not_found_in_gallery").
				WithMeta(map[string]any{"product_id": productID, "media_id": mediaID})
		}
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	filterReset := bson.M{"_id": productID}
	updateReset := bson.M{"$set": bson.M{"gallery.$[].is_primary": false}}
	_, err = d.collection.UpdateOne(ctx, filterReset, updateReset)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	filterSetPrimary := bson.M{"_id": productID, "gallery.media_id": mediaID}
	updateSetPrimary := bson.M{"$set": bson.M{"gallery.$.is_primary": true}}
	_, err = d.collection.UpdateOne(ctx, filterSetPrimary, updateSetPrimary)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}

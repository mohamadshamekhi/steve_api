package mongodbbanner

import (
	"context"
	"errors"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) Create(ctx context.Context, b entity.Banner) error {
	const op = "mongodbbanner.Create"
	b.BeforeCreate()
	_, err := d.collection.InsertOne(ctx, b)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (d DB) Update(ctx context.Context, b *entity.Banner) error {
	const op = "mongodbanner.Update"

	b.UpdatedAt = time.Now()
	filter := bson.M{"_id": b.ID}

	result, err := d.collection.ReplaceOne(ctx, filter, b)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
	}

	return nil
}

func (d DB) FindByID(ctx context.Context, id string) (entity.Banner, error) {
	const op = "mongodbbrand.FindByID"

	var b entity.Banner
	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}

	err := d.collection.FindOne(ctx, filter).Decode(&b)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return b, richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return b, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return b, nil
}

func (d DB) Delete(ctx context.Context, id string) error {
	const op = "mongodbbrand.Delete"

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now(), "is_active": false}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("banner.not_found")
	}

	return nil
}
func (d DB) List(ctx context.Context, isActive *bool) ([]entity.Banner, error) {
	const op = "mongodbbrand.List"

	filter := bson.M{"deleted_at": bson.M{"$exists": false}}

	if isActive != nil {
		filter["is_active"] = *isActive
	}

	opts := options.Find().SetSort(bson.D{{Key: "sort_order", Value: 1}})

	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var banners []entity.Banner
	if err := cursor.All(ctx, &banners); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return banners, nil
}

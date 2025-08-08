package mongodbbrand

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) CheckBrandNameIsExist(ctx context.Context, brandName string) (bool, error) {
	const op = "mongodbbrand.CheckBrandNameIsExist"
	filter := bson.M{"name": brandName}
	var brand entity.Brand

	if err := d.collection.FindOne(ctx, filter).Decode(&brand); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, richerror.New(op).WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"brand_name": brandName})
	}
	return true, nil
}

func (d DB) CheckBrandIDIsExist(ctx context.Context, id string) (bool, error) {
	const op = "mongodbbrand.CheckBrandIDIsExist"
	filter := bson.M{"_id": id}
	var brand entity.Brand

	if err := d.collection.FindOne(ctx, filter).Decode(&brand); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, richerror.New(op).WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}
	return true, nil
}

package mongodbcategory

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	categoryservice "github.com/abolfazlnorzad/steve/service/category"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) All(ctx context.Context, f commonparam.FilterRequest, p commonparam.PaginationRequest) (categoryservice.ListResult, error) {
	const op = "mongodbcategory.All"
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}
	for _, f := range f.Filter {
		for key, value := range f {
			if key == "name" {
				filter[key] = bson.M{"$regex": value, "$options": "i"}
			} else {
				filter[key] = value
			}
		}
	}
	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return categoryservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	offset := p.GetOffset()
	limit := p.GetPageSize()
	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return categoryservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)
	var categories []entity.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return categoryservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return categoryservice.ListResult{
		Categories: categories,
		TotalCount: totalCount,
	}, nil
}

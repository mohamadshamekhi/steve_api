package mongodbbrand

import (
	"context"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	brandservice "github.com/abolfazlnorzad/steve/service/brand"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) All(ctx context.Context, f commonparam.FilterRequest, p commonparam.PaginationRequest) (brandservice.ListResult, error) {
	const op = "brandservice.All"
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}

	for _, item := range f.Filter {
		for key, value := range item {

			switch key {

			case "name":
				if v, ok := value.(string); ok && v != "" {
					filter["name"] = bson.M{
						"$regex":   v,
						"$options": "i",
					}
				}

			case "is_favorite":
				if v, ok := value.(bool); ok {
					filter["$or"] = []bson.M{
						{"is_favorite": v},
						{"is_favorite": bson.M{"$exists": v}},
					}
				}
			}
		}
	}

	fmt.Println("qqq", filter)

	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return brandservice.ListResult{}, richerror.New(op).WithWarpError(err).
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
		return brandservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)
	var brands []entity.Brand
	if err := cursor.All(ctx, &brands); err != nil {
		return brandservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return brandservice.ListResult{
		Brands:     brands,
		TotalCount: totalCount,
	}, nil
}

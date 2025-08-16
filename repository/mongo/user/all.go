package mongodbuser

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	userservice "github.com/abolfazlnorzad/steve/service/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) GetAllUsers(ctx context.Context, f commonparam.FilterRequest, p commonparam.PaginationRequest) (userservice.ListResult, error) {
	const op = "userservice.GetAllUsers"
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}
	for _, filterItem := range f.Filter {
		for key, value := range filterItem {
			if value != nil {
				valueStr := fmt.Sprintf("%v", value)
				if valueStr != "" {
					if key == "name" || key == "phone_number" {
						filter[key] = bson.M{"$regex": valueStr, "$options": "i"}
					} else {
						filter[key] = value
					}
				}
			}
		}
	}
	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return userservice.ListResult{}, richerror.New(op).WithWarpError(err).
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
		return userservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)
	var users []entity.User
	if err := cursor.All(ctx, &users); err != nil {
		return userservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return userservice.ListResult{
		Users:     users,
		TotalCount: totalCount,
	}, nil
}

package mongodbcomment

import (
	"context"
	"fmt"
	"strconv"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	commentservice "github.com/abolfazlnorzad/steve/service/comment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) List(ctx context.Context, f commonparam.FilterRequest, p commonparam.PaginationRequest) (commentservice.ListResult, error) {
	const op = "mongodbcomment.List"

	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}
	for _, filterItem := range f.Filter {
		for key, value := range filterItem {

			valueStr := fmt.Sprintf("%v", value)

			if valueStr == "" {
				continue
			}

			switch key {

			case "is_approved":
				if boolValue, err := strconv.ParseBool(valueStr); err == nil {
					filter[key] = boolValue
				}

			default:
				filter[key] = value
			}
		}
	}

	offset := p.GetOffset()
	limit := p.GetPageSize()
	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	fmt.Println("filter", filter)
	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return commentservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return commentservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var comments []entity.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return commentservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return commentservice.ListResult{
		Comments:   comments,
		TotalCount: uint(totalCount),
	}, nil
}

package mongodbdiscount

import (
	"context"
	"fmt"
	"strconv"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	discountservice "github.com/abolfazlnorzad/steve/service/discount"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) All(ctx context.Context, f commonparam.FilterRequest, p commonparam.PaginationRequest) (discountservice.ListResult, error) {
	const op = "mongodbdiscount.All"

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
			case "code", "minimum_purchase_amount":
				filter[key] = bson.M{"$regex": valueStr, "$options": "i"}

			case "is_active":
				if boolValue, err := strconv.ParseBool(valueStr); err == nil {
					filter[key] = boolValue
				}

			case "type", "value", "starts_at", "expires_at":
				filter[key] = valueStr

			default:
				filter[key] = value
			}
		}
	}

	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return discountservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	offset := p.GetOffset()
	limit := p.GetPageSize()
	projection := bson.M{
		"_id":               1,
		"code":              1,
		"type":              1,
		"value":             1,
		"starts_at":         1,
		"expires_at":        1,
		"total_usage_limit": 1,
		"usage_count":       1,
	}

	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetProjection(projection)

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return discountservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)
	var discounts []discountparam.DiscountsList
	if err := cursor.All(ctx, &discounts); err != nil {
		return discountservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	return discountservice.ListResult{
		Discounts:  discounts,
		TotalCount: totalCount,
	}, nil
}

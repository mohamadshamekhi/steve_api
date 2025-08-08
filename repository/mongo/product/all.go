package mongodbproduct

import (
	"context"
	"fmt"
	"strconv"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	productservice "github.com/abolfazlnorzad/steve/service/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) All(ctx context.Context, f commonparam.FilterRequest, p commonparam.PaginationRequest) (productservice.ListResult, error) {
	const op = "mongodbproduct.All"
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
			case "name", "slug", "en_title":
				filter[key] = bson.M{"$regex": valueStr, "$options": "i"}

			case "is_published", "is_second_hand":
				if boolValue, err := strconv.ParseBool(valueStr); err == nil {
					filter[key] = boolValue
				}

			case "brand_id", "category_id":
				// برای شناسه‌ها، از برابری دقیق استفاده می‌کنیم.
				filter[key] = valueStr

			// می‌توانید case های دیگری برای فیلدهای عددی اضافه کنید.
			// case "min_price":
			// 	if price, err := strconv.Atoi(valueStr); err == nil {
			// 		filter["price"] = bson.M{"$gte": price}
			// 	}

			default:
				// برای هر کلید ناشناخته دیگری، از برابری دقیق استفاده می‌کنیم.
				filter[key] = value
			}
		}
	}

	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return productservice.ListResult{}, richerror.New(op).WithWarpError(err).
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
		return productservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)
	var products []entity.Product
	if err := cursor.All(ctx, &products); err != nil {
		return productservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	return productservice.ListResult{
		Products:   products,
		TotalCount: totalCount,
	}, nil
}

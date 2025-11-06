package mongodbproduct

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	productparam "github.com/abolfazlnorzad/steve/param/product"
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

func (d DB) GetFilteredProducts(ctx context.Context, filter *productparam.StorefrontFilterProduct, p commonparam.PaginationRequest) (productservice.ListResult, error) {
	const op = "mongodbproduct.GetFilteredProducts"

	// ۱. گرفتن پایپ‌لاین پایه (که حالا خروجی‌اش درست است)
	basePipeline, err := d.buildBasePipeline(*filter)
	if err != nil {
		return productservice.ListResult{}, err
	}

	fullPipeline := basePipeline

	sortStage := bson.D{}
	if filter.OrderBy != nil {
		switch *filter.OrderBy {
		case "newest":
			sortStage = bson.D{{"doc.created_at", -1}}
		case "oldest":
			sortStage = bson.D{{"doc.created_at", 1}}
		case "most_expensive":
			sortStage = bson.D{{"min_price_variant", -1}}
		case "least_expensive":
			sortStage = bson.D{{"min_price_variant", 1}}
		case "most_sold":
			sortStage = bson.D{{"doc.sales_count", -1}}
		default:
			sortStage = bson.D{{"doc.created_at", -1}}
		}
	} else {
		sortStage = bson.D{{"doc.created_at", -1}}
	}
	if len(sortStage) > 0 {
		fullPipeline = append(fullPipeline, bson.D{{"$sort", sortStage}})
	}

	offset := p.GetOffset()
	limit := p.GetPageSize()
	facetPipeline := bson.D{
		{"$facet", bson.M{
			"count": bson.A{
				bson.D{{"$count", "total_count"}},
			},
			"data": bson.A{
				bson.D{{"$skip", offset}},
				bson.D{{"$limit", limit}},

				bson.D{{"$replaceRoot", bson.D{
					{"newRoot", bson.D{
						{"$mergeObjects", bson.A{
							"$doc",
							bson.M{
								"_id":      "$_id",
								"variants": "$filtered_variants",
							},
						}},
					}},
				}}},
			},
		}},
	}
	fullPipeline = append(fullPipeline, facetPipeline)

	// ۴. اجرای پایپ‌لاین
	cursor, err := d.collection.Aggregate(ctx, fullPipeline)
	if err != nil {
		return productservice.ListResult{}, richerror.New(op).WithWarpError(err)
	}
	defer cursor.Close(ctx)

	// ۵. دیکد کردن نتایج (بدون تغییر)
	var facetResults []struct {
		Count []struct {
			TotalCount int64 `bson:"total_count"`
		} `bson:"count"`
		Data []entity.Product `bson:"data"`
	}

	if err := cursor.All(ctx, &facetResults); err != nil {
		return productservice.ListResult{}, richerror.New(op).WithWarpError(err)
	}

	// ممکن است هیچ نتیجه‌ای وجود نداشته باشد
	if len(facetResults) == 0 {
		return productservice.ListResult{Products: []entity.Product{}, TotalCount: 0}, nil
	}

	response := productservice.ListResult{
		Products: facetResults[0].Data,
	}
	if len(facetResults[0].Count) > 0 {
		response.TotalCount = facetResults[0].Count[0].TotalCount
	}

	return response, nil
}

func (d DB) buildBasePipeline(filter productparam.StorefrontFilterProduct) ([]bson.D, error) {
	pipeline := []bson.D{}

	// ========== مرحله ۱: Match اولیه (فیلدهای سطح بالای محصول) ==========
	initialMatch := bson.D{{"is_published", true}}
	if filter.BrandID != nil && *filter.BrandID != "" {
		initialMatch = append(initialMatch, bson.E{"brand_id", *filter.BrandID})
	}
	if filter.CategoryID != nil && *filter.CategoryID != "" {
		initialMatch = append(initialMatch, bson.E{"category_id", *filter.CategoryID})
	}
	pipeline = append(pipeline, bson.D{{"$match", initialMatch}})

	// ========== مرحله ۲: فیلتر کردن بر اساس Variants ==========
	pipeline = append(pipeline, bson.D{{"$unwind", "$variants"}})

	variantMatch := bson.D{}
	priceMatch := bson.A{}
	if filter.MinPrice != nil {
		priceMatch = append(priceMatch, bson.D{{"variants.price", bson.D{{"$gte", *filter.MinPrice}}}})
	}
	if filter.MaxPrice != nil {
		priceMatch = append(priceMatch, bson.D{{"variants.price", bson.D{{"$lte", *filter.MaxPrice}}}})
	}
	if len(priceMatch) > 0 {
		variantMatch = append(variantMatch, bson.E{"$and", priceMatch})
	}
	if filter.IsExist != nil {
		if *filter.IsExist {
			variantMatch = append(variantMatch, bson.E{"variants.quantity", bson.D{{"$gt", 0}}})
		} else {
			variantMatch = append(variantMatch, bson.E{"variants.quantity", 0})
		}
	}
	if filter.Color != nil && *filter.Color != "" {
		variantMatch = append(variantMatch, bson.E{"variants.variant_attribute", bson.D{
			{"$elemMatch", bson.D{{"value", *filter.Color}}},
		}})
	}
	if len(variantMatch) > 0 {
		pipeline = append(pipeline, bson.D{{"$match", variantMatch}})
	}

	// ========== مرحله ۳: فیلتر کردن بر اساس Specifications ==========
	if len(filter.Filter) > 0 {
		specMatches := bson.A{}
		for key, val := range filter.Filter {
			specMatches = append(specMatches, bson.D{
				{"specifications", bson.D{
					{"$elemMatch", bson.D{{"key", key}, {"value", val}}},
				}},
			})
		}
		pipeline = append(pipeline, bson.D{{"$match", bson.D{{"$and", specMatches}}}})
	}

	// ========== مرحله ۴: Group (جمع‌آوری مجدد محصول) ==========
	pipeline = append(pipeline, bson.D{{"$group", bson.D{
		{"_id", "$_id"},
		{"filtered_variants", bson.D{{"$push", "$variants"}}},
		{"min_price_variant", bson.D{{"$min", "$variants.price"}}},
		{"doc", bson.D{{"$first", bson.M{
			"name":           "$name",
			"en_title":       "$en_title",
			"slug":           "$slug",
			"brand_id":       "$brand_id",
			"category_id":    "$category_id",
			"gallery":        "$gallery",
			"is_published":   "$is_published",
			"is_second_hand": "$is_second_hand",
			"description":    "$description",
			"specifications": "$specifications",
			"rating":         "$rating",
			"sales_count":    "$sales_count",
			"created_at":     "$created_at",
			"updated_at":     "$updated_at",
			"deleted_at":     "$deleted_at",
		}}}},
	}}})

	// ========== مرحله ۵: Lookup برای پیدا کردن تخفیف‌های مرتبط ==========
	now := time.Now()
	pipeline = append(pipeline, bson.D{
		{"$lookup", bson.M{
			"from": "discounts", // نام کالکشن تخفیف‌های شما
			"let": bson.M{
				"product_id":  "$_id",
				"category_id": "$doc.category_id",
			},
			"pipeline": bson.A{
				bson.M{"$match": bson.M{
					"$expr": bson.M{
						"$and": bson.A{
							// شرط ۱: تخفیف فعال، در بازه زمانی معتبر و دارای ظرفیت استفاده باشد
							bson.M{"$eq": bson.A{"$is_active", true}},
							bson.M{"$lte": bson.A{"$starts_at", now}},
							bson.M{"$gte": bson.A{"$expires_at", now}},
							bson.M{"$gt": bson.A{"$total_usage_limit", "$usage_count"}},

							// شرط ۲: تخفیف قابل اعمال بر روی این محصول باشد (عمومی یا خاص)
							bson.M{"$or": bson.A{
								// حالت الف: تخفیف عمومی است (به هیچ محصول یا دسته‌بندی خاصی محدود نشده)
								bson.M{"$and": bson.A{
									bson.M{"$eq": bson.A{bson.M{"$ifNull": bson.A{"$applicable_product_ids", bson.A{}}}, bson.A{}}},
									bson.M{"$eq": bson.A{bson.M{"$ifNull": bson.A{"$applicable_category_ids", bson.A{}}}, bson.A{}}},
								}},
								// حالت ب: ID محصول در لیست محصولات مجاز تخفیف قرار دارد
								bson.M{"$in": bson.A{"$$product_id", bson.M{"$ifNull": bson.A{"$applicable_product_ids", bson.A{}}}}},
								// حالت ج: ID دسته‌بندی محصول در لیست دسته‌بندی‌های مجاز تخفیف قرار دارد
								bson.M{"$in": bson.A{"$$category_id", bson.M{"$ifNull": bson.A{"$applicable_category_ids", bson.A{}}}}},
							}},
						},
					},
				}},
			},
			"as": "discounts", // نتیجه در فیلد 'discounts' ذخیره می‌شود
		}},
	})

	return pipeline, nil
}

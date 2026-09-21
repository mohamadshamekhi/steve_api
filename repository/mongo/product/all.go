package mongodbproduct

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	productservice "github.com/abolfazlnorzad/steve/service/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// buildListFilter turns the generic key/value filter (as used by the backoffice product
// list and export endpoints) into a Mongo query, always excluding soft-deleted products.
func buildListFilter(f commonparam.FilterRequest) bson.M {
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}
	// Filters that need their own $or (the backward-compatible "field missing
	// means the old default" shapes) are collected here so several of them can
	// coexist instead of overwriting one another's $or.
	andConditions := make([]bson.M, 0, 2)

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
				filter[key] = valueStr

			case "kind":
				// Documents written before the kind field existed are goods;
				// asking for "physical" must therefore also match a missing
				// field, not just an explicit one.
				if valueStr == string(entity.ProductKindPhysical) {
					andConditions = append(andConditions, bson.M{"$or": []bson.M{
						{"kind": bson.M{"$exists": false}},
						{"kind": entity.ProductKindPhysical},
					}})
					continue
				}
				filter[key] = valueStr

			case "sales_channel":
				// Same backward-compatible shape: a document with no channel
				// list was sellable on both channels.
				andConditions = append(andConditions, bson.M{"$or": []bson.M{
					{"sales_channels": bson.M{"$exists": false}},
					{"sales_channels": bson.M{"$size": 0}},
					{"sales_channels": valueStr},
				}})

			default:
				filter[key] = value
			}
		}
	}

	if len(andConditions) > 0 {
		filter["$and"] = andConditions
	}
	return filter
}

func (d DB) All(ctx context.Context, f commonparam.FilterRequest, p commonparam.PaginationRequest) (productservice.ListResult, error) {
	const op = "mongodbproduct.All"
	filter := buildListFilter(f)

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

func (d DB) AllForExport(ctx context.Context, f commonparam.FilterRequest) ([]entity.Product, error) {
	const op = "mongodbproduct.AllForExport"
	filter := buildListFilter(f)
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var products []entity.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	return products, nil
}

func (d DB) GetFilteredProducts(ctx context.Context, filter *productparam.StorefrontFilterProduct, p commonparam.PaginationRequest) (productservice.ListResult, error) {
	const op = "mongodbproduct.GetFilteredProducts"

	pipeline, err := d.buildBasePipeline(*filter)
	if err != nil {
		return productservice.ListResult{}, err
	}

	sortStage := bson.D{}
	if filter.OrderBy != nil {
		switch *filter.OrderBy {
		case "newest":
			sortStage = bson.D{{"created_at", -1}, {"_id", -1}}
		case "oldest":
			sortStage = bson.D{{"created_at", 1}, {"_id", 1}}
		case "most_expensive":
			// مرتب‌سازی بر اساس گران‌ترین واریانت محصول
			sortStage = bson.D{{"calculated_max_price", -1}, {"created_at", -1}, {"_id", -1}}
		case "least_expensive":
			// مرتب‌سازی بر اساس ارزان‌ترین واریانت محصول
			sortStage = bson.D{{"calculated_min_price", 1}, {"created_at", -1}, {"_id", -1}}
		case "most_sold":
			sortStage = bson.D{{"sales_count", -1}, {"created_at", -1}, {"_id", -1}}
		default:
			sortStage = bson.D{{"created_at", -1}, {"_id", -1}}
		}
	} else {
		sortStage = bson.D{{"created_at", -1}, {"_id", -1}}
	}

	if len(sortStage) > 0 {
		pipeline = append(pipeline, bson.D{{"$sort", sortStage}})
	}

	// ۳. افزودن مرحله صفحه‌بندی و Facet
	offset := p.GetOffset()
	limit := p.GetPageSize()

	facetStage := bson.D{
		{"$facet", bson.M{
			"count": bson.A{
				bson.D{{"$count", "total_count"}},
			},
			"data": bson.A{
				bson.D{{"$skip", offset}},
				bson.D{{"$limit", limit}},
			},
		}},
	}
	pipeline = append(pipeline, facetStage)

	// ۴. اجرای پایپ‌لاین
	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return productservice.ListResult{}, richerror.New(op).WithWarpError(err)
	}
	defer cursor.Close(ctx)

	// ۵. دیکود کردن نتایج
	var facetResults []struct {
		Count []struct {
			TotalCount int64 `bson:"total_count"`
		} `bson:"count"`
		Data []entity.Product `bson:"data"`
	}

	if err := cursor.All(ctx, &facetResults); err != nil {
		return productservice.ListResult{}, richerror.New(op).WithWarpError(err)
	}

	if len(facetResults) == 0 {
		return productservice.ListResult{Products: []entity.Product{}, TotalCount: 0}, nil
	}

	totalCount := int64(0)
	if len(facetResults[0].Count) > 0 {
		totalCount = facetResults[0].Count[0].TotalCount
	}

	return productservice.ListResult{
		Products:   facetResults[0].Data,
		TotalCount: totalCount,
	}, nil
}

func (d DB) buildBasePipeline(filter productparam.StorefrontFilterProduct) ([]bson.D, error) {
	pipeline := []bson.D{}

	pipeline = append(pipeline, bson.D{
		{"$addFields", bson.M{
			"calculated_min_price": bson.M{"$min": "$variants.price"},
			"calculated_max_price": bson.M{"$max": "$variants.price"},
			"sales_count":          bson.M{"$ifNull": bson.A{"$sales_count", 0}},
		}},
	})

	initialMatch := bson.D{
		{"is_published", true},
		{"deleted_at", bson.M{"$exists": false}},
		{"$or", websiteChannelCondition()["$or"]},
	}
	if filter.BrandID != nil && *filter.BrandID != "" {
		initialMatch = append(initialMatch, bson.E{"brand_id", *filter.BrandID})
	}
	if filter.CategoryID != nil && *filter.CategoryID != "" {
		initialMatch = append(initialMatch, bson.E{"category_id", *filter.CategoryID})
	}
	if filter.IsSecondHand != nil {
		initialMatch = append(initialMatch, bson.E{"is_second_hand", *filter.IsSecondHand})
	}

	if len(filter.Filter) > 0 {
		specMatches := bson.A{}
		for key, val := range filter.Filter {
			specMatches = append(specMatches, bson.D{
				{"specifications", bson.D{
					{"$elemMatch", bson.D{{"key", key}, {"value", val}}},
				}},
			})
		}
		initialMatch = append(initialMatch, bson.E{"$and", specMatches})
	}

	pipeline = append(pipeline, bson.D{{"$match", initialMatch}})

	hasVariantFilter := (filter.MinPrice != nil) ||
		(filter.MaxPrice != nil) ||
		(filter.IsExist != nil) ||
		(filter.Color != nil && *filter.Color != "")

	if hasVariantFilter {
		pipeline = append(pipeline, bson.D{{"$unwind", "$variants"}})

		variantMatch := bson.D{}

		if filter.MinPrice != nil {
			variantMatch = append(variantMatch, bson.E{"variants.price", bson.D{{"$gte", *filter.MinPrice}}})
		}
		if filter.MaxPrice != nil {
			variantMatch = append(variantMatch, bson.E{"variants.price", bson.D{{"$lte", *filter.MaxPrice}}})
		}

		if filter.IsExist != nil {
			if *filter.IsExist {
				variantMatch = append(variantMatch, bson.E{"variants.quantity", bson.D{{"$gt", 0}}})
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

		pipeline = append(pipeline, bson.D{{"$group", bson.D{
			{"_id", "$_id"},
			{"filtered_variants", bson.D{{"$push", "$variants"}}},

			{"calculated_min_price", bson.D{{"$first", "$calculated_min_price"}}},
			{"calculated_max_price", bson.D{{"$first", "$calculated_max_price"}}},
			{"sales_count", bson.D{{"$first", "$sales_count"}}},
			{"created_at", bson.D{{"$first", "$created_at"}}},

			{"doc_root", bson.D{{"$first", "$$ROOT"}}},
		}}})

		pipeline = append(pipeline, bson.D{
			{"$replaceRoot", bson.D{
				{"newRoot", bson.D{
					{"$mergeObjects", bson.A{
						"$doc_root",
						bson.M{
							"variants":             "$filtered_variants",
							"calculated_min_price": "$calculated_min_price",
							"calculated_max_price": "$calculated_max_price",
						},
					}},
				}},
			}},
		})
	}

	now := time.Now()
	pipeline = append(pipeline, bson.D{
		{"$lookup", bson.M{
			"from": "discounts",
			"let": bson.M{
				"product_id":  "$_id",
				"category_id": "$category_id",
			},
			"pipeline": bson.A{
				bson.M{"$match": bson.M{
					"$expr": bson.M{
						"$and": bson.A{
							bson.M{"$eq": bson.A{"$is_active", true}},
							bson.M{"$lte": bson.A{"$starts_at", now}},
							bson.M{"$gte": bson.A{"$expires_at", now}},
							bson.M{"$gt": bson.A{"$total_usage_limit", "$usage_count"}},

							bson.M{"$or": bson.A{
								bson.M{"$and": bson.A{
									bson.M{"$eq": bson.A{bson.M{"$ifNull": bson.A{"$applicable_product_ids", bson.A{}}}, bson.A{}}},
									bson.M{"$eq": bson.A{bson.M{"$ifNull": bson.A{"$applicable_category_ids", bson.A{}}}, bson.A{}}},
								}},
								bson.M{"$in": bson.A{"$$product_id", bson.M{"$ifNull": bson.A{"$applicable_product_ids", bson.A{}}}}},
								bson.M{"$in": bson.A{"$$category_id", bson.M{"$ifNull": bson.A{"$applicable_category_ids", bson.A{}}}}},
							}},
						},
					},
				}},
			},
			"as": "discounts",
		}},
	})

	return pipeline, nil
}

func (d DB) GetSearchProduct(ctx context.Context, search commonparam.SearchRequest) (productservice.ListResult, error) {
	const op = "mongodbproduct.GetSearchProduct"

	pipeline := mongo.Pipeline{}

	match := bson.M{
		"is_published": true,
	}

	andConditions := bson.A{
		websiteChannelCondition(),
		bson.M{
			"$or": bson.A{
				bson.M{"deleted_at": bson.M{"$exists": false}},
				bson.M{"deleted_at": nil},
			},
		},
	}

	if search.Search != "" {
		q := strings.TrimSpace(search.Search)
		andConditions = append(andConditions, bson.M{
			"$or": bson.A{
				bson.M{"name": bson.M{"$regex": q, "$options": "i"}},
				bson.M{"en_title": bson.M{"$regex": q, "$options": "i"}},
				bson.M{"slug": bson.M{"$regex": q, "$options": "i"}},
			},
		})
	}

	if len(andConditions) > 0 {
		match["$and"] = andConditions
	}

	pipeline = append(pipeline, bson.D{{"$match", match}})

	pipeline = append(pipeline, bson.D{
		{"$sort", bson.D{{"created_at", -1}, {"_id", -1}}},
	})

	pipeline = append(pipeline, bson.D{
		{"$limit", 10},
	})

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return productservice.ListResult{}, richerror.New(op).WithWarpError(err)
	}
	defer cursor.Close(ctx)

	var products []entity.Product
	if err := cursor.All(ctx, &products); err != nil {
		return productservice.ListResult{}, richerror.New(op).WithWarpError(err)
	}

	return productservice.ListResult{
		Products: products,
	}, nil
}

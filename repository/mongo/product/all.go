package mongodbproduct

import (
	"context"
	"fmt"
	"strconv"

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

	basePipeline, err := d.buildBasePipeline(*filter)
	if err != nil {
		return productservice.ListResult{}, err
	}

	offset := p.GetOffset()
	limit := p.GetPageSize()
	facetPipeline := []bson.D{
		bson.D{ // این wrapper خطای کامپایلر را رفع می‌کند.
			{"$facet", bson.M{ // از bson.M برای آرگومان‌های داخلی $facet استفاده شده است
				// Pipeline برای شمارش کل نتایج
				"count": bson.A{
					bson.D{{"$count", "total_count"}},
				},
				// Pipeline برای داده‌های صفحه‌بندی شده
				"data": bson.A{
					bson.D{{"$skip", offset}},
					bson.D{{"$limit", limit}},
					// جایگزینی ریشه با داکیومنت اصلی محصول
					bson.D{{"$replaceRoot", bson.D{{"newRoot", "$doc"}}}},
				},
			}},
		},
	}

	fullPipeline := append(basePipeline, facetPipeline...)

	cursor, err := d.collection.Aggregate(ctx, fullPipeline)
	if err != nil {
		return productservice.ListResult{}, richerror.New(op).WithWarpError(err)
	}
	defer cursor.Close(ctx)

	// ۴. دیکد کردن نتایج $facet
	var facetResults []struct {
		Count []struct {
			TotalCount int64 `bson:"total_count"`
		} `bson:"count"`
		Data []entity.Product `bson:"data"`
	}

	if err := cursor.All(ctx, &facetResults); err != nil {
		return productservice.ListResult{}, richerror.New(op).WithWarpError(err)
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

	// ۱. Match اولیه (فیلدهای سطح بالا)
	initialMatch := bson.D{{"is_published", true}}

	if filter.BrandID != nil && *filter.BrandID != "" {
		initialMatch = append(initialMatch, bson.E{"brand_id", *filter.BrandID})
	}
	if filter.CategoryID != nil && *filter.CategoryID != "" {
		initialMatch = append(initialMatch, bson.E{"category_id", *filter.CategoryID})
	}

	pipeline = append(pipeline, bson.D{{"$match", initialMatch}})

	// ۲. فاز Variants (Unwind و Match)
	pipeline = append(pipeline, bson.D{{"$unwind", "$variants"}})

	variantMatch := bson.D{}

	// فیلتر قیمت
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

	// فیلتر موجودی
	if filter.IsExist != nil {
		if *filter.IsExist {
			variantMatch = append(variantMatch, bson.E{"variants.quantity", bson.D{{"$gt", 0}}})
		} else {
			variantMatch = append(variantMatch, bson.E{"variants.quantity", 0})
		}
	}

	// --- FIX: فیلتر رنگ در variant_attribute ---
	if filter.Color != nil && *filter.Color != "" {
		// از $elemMatch استفاده می‌کنیم تا درون آرایه variants.variant_attribute به دنبال مقدار رنگ بگردیم.
		variantMatch = append(variantMatch, bson.E{"variants.variant_attribute", bson.D{
			{"$elemMatch", bson.D{
				// فرض می‌کنیم مقدار رنگ مستقیماً در فیلد "value" درون یکی از المنت‌های variant_attribute قرار دارد.
				// (مثلاً {key: "color", value: "red"} یا {key: "128gig", value: "red"})
				{"value", *filter.Color},
			}},
		}})
	}
	// --- END FIX ---

	if len(variantMatch) > 0 {
		pipeline = append(pipeline, bson.D{{"$match", variantMatch}})
	}

	// ۳. فاز فیلترینگ داینامیک Specifications
	if len(filter.Filter) > 0 {
		specMatches := bson.A{}
		for key, val := range filter.Filter {
			specMatches = append(specMatches, bson.D{
				{"specifications", bson.D{
					{"$elemMatch", bson.D{
						{"key", key},
						{"value", val},
					}},
				}},
			})
		}

		if len(specMatches) > 0 {
			pipeline = append(pipeline, bson.D{{"$match", bson.D{{"$and", specMatches}}}})
		}
	}

	// --- ۴. فاز Group (اصلاح شده) ---
	// مشکل دیکدینگ را با بازسازی دستی داکیومنت در اینجا حل می‌کنیم تا فیلد variants خراب شده وارد "doc" نشود.
	pipeline = append(pipeline, bson.D{{"$group", bson.D{
		{"_id", "$_id"},
		// محاسبه کمترین قیمت ورینت‌های فیلترشده (برای مرتب‌سازی)
		{"min_price_variant", bson.D{{"$min", "$variants.price"}}},
		// نگه‌داشتن فیلدهای مورد نیاز محصول به صورت مجزا به جای استفاده از $$ROOT
		{"doc", bson.D{{"$first", bson.M{
			"_id":            "$_id",
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
			"created_at":     "$created_at",
			"updated_at":     "$updated_at",
			"deleted_at":     "$deleted_at",
			// مهم: فیلد "variants" را عمداً از اینجا حذف می‌کنیم
		}}},
		}}}})

	// --- ۵. فاز Sort (مرتب‌سازی نهایی) ---
	sortStage := bson.D{}
	if filter.OrderBy != nil {
		if filter.OrderBy != nil {
			switch *filter.OrderBy {
			case "newest":
				sortStage = bson.D{{"doc.created_at", -1}}
			case "oldest":
				sortStage = bson.D{{"doc.created_at", 1}}
			case "most_expensive":
				sortStage = bson.D{{"min_price_variant", -1}} // گران‌ترین
			case "least_expensive":
				sortStage = bson.D{{"min_price_variant", 1}} // ارزان‌ترین
			default:
				sortStage = bson.D{{"doc.created_at", -1}}
			}
		} else {
			sortStage = bson.D{{"doc.created_at", -1}}
		}
		pipeline = append(pipeline, bson.D{{"$sort", sortStage}})

		return pipeline, nil
	}
	return pipeline, nil
}

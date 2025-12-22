package mongodbproduct

import (
	"context"
	"errors"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Product, error) {
	const op = "mongodbproduct.FindByID"
	filter := bson.M{
		"$or": []bson.M{
			{"_id": id},
		},
		"deleted_at": bson.M{"$exists": false},
	}
	var p entity.Product

	if err := d.collection.FindOne(ctx, filter).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return p, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}

	return p, nil
}

func (d DB) FindBySlug(ctx context.Context, slug string) (entity.Product, error) {
	const op = "mongodbproduct.FindBySlug"
	filter := bson.M{
		"$or": []bson.M{
			{"slug": slug},
		},
		"deleted_at": bson.M{"$exists": false},
	}
	var p entity.Product

	if err := d.collection.FindOne(ctx, filter).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return p, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"slug": slug})
	}

	return p, nil
}

func (d DB) FindManyByIDs(ctx context.Context, ids []string) ([]entity.Product, error) {
	const op = "mongodbproduct.FindManyByIDs"
	if len(ids) == 0 {
		return []entity.Product{}, nil
	}

	filter := bson.M{
		"_id":        bson.M{"$in": ids},
		"deleted_at": bson.M{"$exists": false},
	}

	cursor, err := d.collection.Find(ctx, filter)
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

func (d DB) FindNonExistingIDs(ctx context.Context, ids []string) ([]string, error) {
	const op = "mongodbproduct.FindNonExistingIDs"

	if len(ids) == 0 {
		return nil, nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	opts := options.Find().SetProjection(bson.M{"_id": 1})

	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID string `bson:"_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	existingIDSet := make(map[string]struct{}, len(results))
	for _, result := range results {
		existingIDSet[result.ID] = struct{}{}
	}

	var notFoundIDs []string
	for _, id := range ids {
		if _, found := existingIDSet[id]; !found {
			notFoundIDs = append(notFoundIDs, id)
		}
	}

	return notFoundIDs, nil

}

func (d DB) GetVariantsByIDs(ctx context.Context, variantIDs []string) (map[string]entity.Variant, error) {
	const op = "mongodbproduct.GetVariantsByIDs"

	if len(variantIDs) == 0 {
		return make(map[string]entity.Variant), nil
	}

	filter := bson.M{"variants.id": bson.M{"$in": variantIDs}}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	resultMap := make(map[string]entity.Variant)

	for cursor.Next(ctx) {
		var product entity.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
		}

		for _, variant := range product.Variants {
			for _, requestedID := range variantIDs {
				if variant.ID == requestedID {
					resultMap[variant.ID] = variant
					break
				}
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return resultMap, nil
}

func (d DB) GetProductCategoriesForCartItems(ctx context.Context, items []entity.CartItem) (map[string][]string, error) {
	const op = "mongodbproduct.GetProductCategoriesForCartItems"

	if len(items) == 0 {
		return make(map[string][]string), nil
	}

	productIDs := make([]string, 0, len(items))
	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
	}

	filter := bson.M{"_id": bson.M{"$in": productIDs}}

	projection := bson.M{"_id": 1, "category_id": 1}
	opts := options.Find().SetProjection(projection)

	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	resultMap := make(map[string][]string)
	for cursor.Next(ctx) {
		var product struct {
			ID         string `bson:"_id"`
			CategoryID string `bson:"category_id"`
		}
		if err := cursor.Decode(&product); err != nil {
			return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
		}
		resultMap[product.ID] = []string{product.CategoryID}
	}

	if err := cursor.Err(); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return resultMap, nil
}

func (d DB) GetProductsByIDs(ctx context.Context, ids []string) ([]entity.Product, error) {
	const op = "mongodbproduct.GetProductsByIDs"
	if len(ids) == 0 {
		return []entity.Product{}, nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err)
	}
	defer cursor.Close(ctx)

	var products []entity.Product

	for cursor.Next(ctx) {
		var product entity.Product

		if err := cursor.Decode(&product); err != nil {
			return nil, richerror.New(op).WithWarpError(err)
		}

		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, richerror.New(op).WithWarpError(err)
	}

	return products, nil
}

func (d DB) FindProductsIDsByCategory(ctx context.Context, categoryID string) (string, error) {
	const op = "mongodbproduct.FindProductsIDsByCategory"
	filter := bson.M{
		"$or": []bson.M{
			{"category_id": categoryID},
		},
		"deleted_at": bson.M{"$exists": false},
	}

	var p entity.Product
	if err := d.collection.FindOne(ctx, filter).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p.ID, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		fmt.Println("err ", err)
		return p.ID, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"category_id": categoryID})
	}

	return p.ID, nil
}

func (d DB) GetFilterableAttributes(ctx context.Context, categoryID string) ([]productparam.AggregatedFilterOption, error) {
	const op = "productrepo.GetFilterableAttributes"

	pipeline := mongo.Pipeline{

		bson.D{
			{"$match", bson.M{
				"$or": []bson.M{
					{"category_id": categoryID},
					{"ancestors._id": categoryID},
				},
				"is_published": true,
				"deleted_at":   bson.M{"$exists": false},
			}},
		},
		bson.D{
			{"$unwind", "$specifications"},
		},

		bson.D{
			{"$group", bson.M{
				"_id":    "$specifications.key",
				"values": bson.M{"$addToSet": "$specifications.value"},
			}},
		},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var results []productparam.AggregatedFilterOption
	if err := cursor.All(ctx, &results); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return results, nil
}

func (d DB) GetFilterMobileListWithMainInfo(ctx context.Context, name *string, brandID *string) ([]entity.Product, error) {
	const op = "productrepo.GetFilterMobileListWithMainInfo"
	filter := bson.M{}

	if name != nil {
		filter["name"] = bson.M{"$regex": *name, "$options": "i"}
	}

	if brandID != nil {
		filter["brand_id"] = *brandID
	}
	filter["deleted_at"] = bson.M{"$exists": false}
	filter["is_published"] = true
	projection := bson.M{
		"id":             1,
		"name":           1,
		"slug":           1,
		"brand_id":       1,
		"category_id":    1,
		"gallery":        1,
		"is_published":   1,
		"is_second_hand": 1,
	}
	fmt.Println("projection", filter)
	cursor, err := d.collection.Find(ctx, filter, &options.FindOptions{
		Projection: projection,
		Sort:       bson.D{{"name", 1}},
	})
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err)
	}
	var products []entity.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	return products, nil
}

func (d DB) GetMinAdnMaxPriceOfProductWithSpecialCategory(ctx context.Context, categoryID string) (uint64, uint64, error) {
	const op = "productrepo.GetMinAdnMaxPriceOfProductWithSpecialCategory"

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"category_id":  categoryID,
			"deleted_at":   bson.M{"$exists": false},
			"is_published": true,
			"variants":     bson.M{"$exists": true, "$ne": bson.A{}},
		}}},
		{{Key: "$unwind", Value: "$variants"}},
		{{Key: "$group", Value: bson.M{
			"_id": nil,
			"min": bson.M{"$min": "$variants.price"},
			"max": bson.M{"$max": "$variants.price"},
		}}},
	}

	cur, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cur.Close(ctx)

	if !cur.Next(ctx) {
		return 0, 0, nil
	}

	var res struct {
		Min int64 `bson:"min"`
		Max int64 `bson:"max"`
	}
	if err := cur.Decode(&res); err != nil {
		return 0, 0, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	return uint64(res.Min), uint64(res.Max), nil
}

func (d DB) GetProductIDsByCategories(ctx context.Context, categoryIDs []string) (map[string][]string, error) {

	filter := bson.M{"category_id": bson.M{"$in": categoryIDs}}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find products by categories: %w", err)
	}
	defer cursor.Close(ctx)

	categoryToProductIDs := make(map[string][]string)

	// پردازش نتایج query
	for cursor.Next(ctx) {
		var product entity.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, fmt.Errorf("failed to decode product: %w", err)
		}

		// بررسی دسته‌بندی محصول
		if product.CategoryID != "" {
			categoryToProductIDs[product.CategoryID] = append(categoryToProductIDs[product.CategoryID], product.ID)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return categoryToProductIDs, nil
}

//// GetPricesOfProducts پیاده‌سازی متد اصلی
//func (d DB) GetPricesOfProducts(ctx context.Context, prodIDSlice []string, variantIDSlice []string) (map[string]map[string]float64, error) {
//
//}

// GetPricesOfProducts قیمت تمام تنوع‌های مربوط به شناسه‌های محصول ورودی را برمی‌گرداند.
//
// این متد شناسه‌های محصول (`prodIDSlice`) را دریافت کرده و با یک کوئری به دیتابیس،
// تمام محصولاتی که شناسه‌ی آن‌ها در این لیست وجود دارد را پیدا می‌کند.
// سپس برای هر محصول، قیمت تمام تنوع‌هایش را استخراج کرده و در ساختار map[productID]map[variantID]price برمی‌گرداند.
//
// نکته: پارامتر `variantIDSlice` در این پیاده‌سازی مستقیماً در کوئری دیتابیس استفاده نمی‌شود،
// زیرا بهینه‌تر است که کل داکیومنت محصول را یکجا بگیریم و سپس در لایه سرویس، تنوع‌های مورد نظر را فیلتر کنیم.
// این کار جلوی کوئری‌های پیچیده و متعدد به دیتابیس را می‌گیرد.
func (d DB) GetPricesOfProducts(ctx context.Context, prodIDSlice []string, variantIDSlice []string) (map[string]map[string]float64, error) {
	const op = "productrepo.GetPricesOfProducts"
	// 1. آماده سازی map برای ذخیره قیمت‌ها
	productPrices := make(map[string]map[string]float64)

	var products []entity.Product
	filter := bson.M{
		"_id": bson.M{
			"$in": prodIDSlice,
		},
	}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to fetch products: %w", op, err)
	}
	if err := cursor.All(ctx, &products); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	// 3. پردازش داده‌ها و استخراج قیمت‌ها
	for _, product := range products {
		// برای هر محصول، ورینت‌های آن را پردازش می‌کنیم
		for _, variant := range product.Variants {
			// اگر ورینت در لیست ورینت‌ها باشد، قیمت آن را استخراج می‌کنیم
			if contains(variantIDSlice, variant.ID) {
				// اگر محصول در نقشه موجود نباشد، آن را اضافه می‌کنیم
				if _, exists := productPrices[product.ID]; !exists {
					productPrices[product.ID] = make(map[string]float64)
				}

				// ذخیره قیمت ورینت در نقشه
				productPrices[product.ID][variant.ID] = float64(variant.Price)
			}
		}
	}

	return productPrices, nil
}

// تابع کمکی برای چک کردن وجود یک ورینت در آرایه
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

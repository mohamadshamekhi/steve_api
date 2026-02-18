package mongodbarticle

import (
	"context"
	"strings"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	articleservice "github.com/abolfazlnorzad/steve/service/article"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) All(ctx context.Context, f commonparam.FilterRequest, p commonparam.PaginationRequest) (articleservice.ListResult, error) {
	const op = "mongodbarticle.All"
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}
	for _, f := range f.Filter {
		for key, value := range f {
			if key == "title" {
				filter[key] = bson.M{"$regex": value, "$options": "i"}
			} else {
				filter[key] = value
			}
		}
	}
	var sortValue any
	if filter["sort"] == "desc" {
		sortValue = -1
	} else {
		sortValue = 1
	}
	delete(filter, "sort")
	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return articleservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	offset := p.GetOffset()
	limit := p.GetPageSize()

	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: sortValue}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return articleservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)
	var articles []entity.Article
	if err := cursor.All(ctx, &articles); err != nil {
		return articleservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	return articleservice.ListResult{
		Articles:   articles,
		TotalCount: totalCount,
	}, nil
}

func (d DB) GetSearchArticles(ctx context.Context, search commonparam.SearchRequest) (articleservice.ListResult, error) {
	const op = "mongodbarticle.GetSearchArticle"

	pipeline := mongo.Pipeline{}

	match := bson.M{
		"is_published": true,
	}

	andConditions := bson.A{
		bson.M{
			"$or": bson.A{
				bson.M{"deleted_at": bson.M{"$exists": false}},
				bson.M{"deleted_at": time.Time{}},
			},
		},
	}

	if search.Search != "" {
		q := strings.TrimSpace(search.Search)
		andConditions = append(andConditions, bson.M{
			"$or": bson.A{
				bson.M{"title": bson.M{"$regex": q, "$options": "i"}},
				bson.M{"slug": bson.M{"$regex": q, "$options": "i"}},
			},
		})
	}

	match["$and"] = andConditions

	pipeline = append(pipeline, bson.D{{"$match", match}})

	pipeline = append(pipeline, bson.D{
		{"$sort", bson.D{{"created_at", -1}, {"_id", -1}}},
	})

	pipeline = append(pipeline, bson.D{
		{"$limit", 10},
	})

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return articleservice.ListResult{}, richerror.New(op).WithWarpError(err)
	}
	defer cursor.Close(ctx)

	var articles []entity.Article
	if err := cursor.All(ctx, &articles); err != nil {
		return articleservice.ListResult{}, richerror.New(op).WithWarpError(err)
	}

	return articleservice.ListResult{
		Articles: articles,
	}, nil
}

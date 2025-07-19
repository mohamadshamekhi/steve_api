package mongodbattribute

import (
	"context"
	"errors"
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	attributeservice "github.com/abolfazlnorzad/steve/service/attribute"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const collectionName = "attributesets"

type DB struct {
	collection *mongo.Collection
}

func New(c *mongodb.DB) *DB {
	collection := c.Connection().Collection(collectionName)

	return &DB{
		collection: collection,
	}
}

func (d DB) Create(ctx context.Context, attribute entity.AttributeSet) error {
	const op = "mongodbattribute.Create"
	_, err := d.collection.InsertOne(ctx, attribute)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"attribute": attribute}).WithKind(richerror.KindUnexpected)
	}
	return nil
}
func (d DB) Update(ctx context.Context, attribute entity.AttributeSet) error {
	const op = "mongodbattribute.Update"
	filter := bson.M{"_id": attribute.ID}
	attribute.UpdatedAt = time.Now()
	result, err := d.collection.ReplaceOne(ctx, filter, attribute)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found")
	}

	return nil
}
func (d DB) All(ctx context.Context, f commonparam.FilterRequest, p commonparam.PaginationRequest) (attributeservice.ListResult, error) {
	const op = "mongodbattribute.All"
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}
	for _, f := range f.Filter {
		for key, value := range f {
			if key == "name" {
				filter[key] = bson.M{"$regex": value, "$options": "i"}
			} else {
				filter[key] = value
			}
		}
	}
	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return attributeservice.ListResult{}, richerror.New(op).WithWarpError(err).
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
		return attributeservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)
	var attrs []entity.AttributeSet
	if err := cursor.All(ctx, &attrs); err != nil {
		return attributeservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return attributeservice.ListResult{
		Attributes: attrs,
		TotalCount: totalCount,
	}, nil
}
func (d DB) Delete(ctx context.Context, id string) error {
	const op = "mongodbattribute.Delete"

	filter := bson.M{"_id": id}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if result.DeletedCount == 0 {
		return richerror.New(op).
			WithMessageKey("media.not_found").
			WithKind(richerror.KindNotFound).
			WithMeta(map[string]any{"id": id})
	}

	return nil
}

func (d DB) FindByID(ctx context.Context, id string) (entity.AttributeSet, error) {
	const op = "mongodbattribute.FindByID"
	filter := bson.M{"_id": id}
	var atte entity.AttributeSet

	if err := d.collection.FindOne(ctx, filter).Decode(&atte); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return atte, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return atte, richerror.New(op).WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}
	return atte, nil
}

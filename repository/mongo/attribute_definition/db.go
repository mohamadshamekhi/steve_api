package mongodbattributedefinition

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
)

const collectionName = "attributes_definitions"

type DB struct {
	collection *mongo.Collection
}

func New(c *mongodb.DB) *DB {
	collection := c.Connection().Collection(collectionName)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "key", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		panic(err)
	}
	return &DB{
		collection: collection,
	}
}

func (d DB) Create(ctx context.Context, attribute entity.AttributeDefinition) error {
	const op = "mongodbattributedefinition.Create"
	_, err := d.collection.InsertOne(ctx, attribute)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"attribute": attribute}).WithKind(richerror.KindUnexpected)
	}
	return nil
}
func (d DB) Update(ctx context.Context, attribute entity.AttributeDefinition) error {
	const op = "mongodbattributedefinition.Update"
	filter := bson.M{"_id": attribute.ID}
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
func (d DB) All(ctx context.Context, f commonparam.FilterRequest, p commonparam.PaginationRequest) (attributeservice.ListDefResult, error) {
	const op = "mongodbattributedefinition.All"
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
		return attributeservice.ListDefResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	offset := p.GetOffset()
	limit := p.GetPageSize()
	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return attributeservice.ListDefResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)
	var attrs []entity.AttributeDefinition
	if err := cursor.All(ctx, &attrs); err != nil {
		return attributeservice.ListDefResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return attributeservice.ListDefResult{
		Attributes: attrs,
		TotalCount: totalCount,
	}, nil
}
func (d DB) Delete(ctx context.Context, id string) error {
	const op = "mongodbattributedefinition.Delete"

	filter := bson.M{"_id": id}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if result.DeletedCount == 0 {
		return richerror.New(op).
			WithMessageKey("global.not-found").
			WithKind(richerror.KindNotFound).
			WithMeta(map[string]any{"id": id})
	}

	return nil

}
func (d DB) FindByIDs(ctx context.Context, ids []string) ([]entity.AttributeDefinition, error) {
	const op = "mongodbattributedefinition.FindByIDs"
	if len(ids) == 0 {
		return []entity.AttributeDefinition{}, nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var fetched []entity.AttributeDefinition
	if err := cursor.All(ctx, &fetched); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	byID := make(map[string]entity.AttributeDefinition, len(fetched))
	for _, a := range fetched {
		byID[a.ID] = a
	}
	attrs := make([]entity.AttributeDefinition, 0, len(ids))
	for _, id := range ids {
		if a, ok := byID[id]; ok {
			attrs = append(attrs, a)
		}
	}

	return attrs, nil
}
func (d DB) FindByKey(ctx context.Context, key string) (entity.AttributeDefinition, error) {
	const op = "mongodbattributedefinition.FindByKey"
	filter := bson.M{"key": key}
	var atte entity.AttributeDefinition

	if err := d.collection.FindOne(ctx, filter).Decode(&atte); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return atte, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return atte, richerror.New(op).WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"key": key})
	}
	return atte, nil
}
func (d DB) FindByID(ctx context.Context, id string) (entity.AttributeDefinition, error) {
	const op = "mongodbattributedefinition.FindByID"
	filter := bson.M{"_id": id}
	var atte entity.AttributeDefinition

	if err := d.collection.FindOne(ctx, filter).Decode(&atte); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return atte, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return atte, richerror.New(op).WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}
	return atte, nil
}

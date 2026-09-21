package mongodbsupplier

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Supplier, error) {
	const op = "mongodbsupplier.FindByID"

	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}

	var s entity.Supplier
	if err := d.collection.FindOne(ctx, filter).Decode(&s); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Supplier{}, richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("supplier.not_found")
		}
		return entity.Supplier{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return s, nil
}

// FindManyByIDs resolves supplier names for list/report responses in one round
// trip. Missing IDs are simply absent from the map.
func (d DB) FindManyByIDs(ctx context.Context, ids []string) (map[string]entity.Supplier, error) {
	const op = "mongodbsupplier.FindManyByIDs"

	result := make(map[string]entity.Supplier, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	cursor, err := d.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var suppliers []entity.Supplier
	if err := cursor.All(ctx, &suppliers); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	for _, s := range suppliers {
		result[s.ID] = s
	}
	return result, nil
}

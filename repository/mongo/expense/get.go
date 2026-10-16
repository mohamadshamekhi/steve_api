package mongodbexpense

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Expense, error) {
	const op = "mongodbexpense.FindByID"

	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}

	var e entity.Expense
	if err := d.collection.FindOne(ctx, filter).Decode(&e); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Expense{}, richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("expense.not_found")
		}
		return entity.Expense{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return e, nil
}

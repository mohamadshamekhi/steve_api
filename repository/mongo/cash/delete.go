package mongodbcash

import (
	"context"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// Delete removes one cashbox line. Only manual lines are deletable; the
// service refuses auto-recorded ones, which belong to their source document.
func (d DB) Delete(ctx context.Context, id string) error {
	const op = "mongodbcash.Delete"

	result, err := d.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if result.DeletedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("cash.not_found")
	}
	return nil
}

// FindByID is used before a delete to check whether the line is auto-recorded.
func (d DB) FindByID(ctx context.Context, id string) (bool, error) {
	const op = "mongodbcash.FindByID"

	var doc struct {
		IsAuto bool `bson:"is_auto"`
	}
	if err := d.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc); err != nil {
		return false, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindNotFound).WithMessageKey("cash.not_found")
	}
	return doc.IsAuto, nil
}

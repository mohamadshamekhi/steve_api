package mongodbsalesreturn

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.SalesReturn, error) {
	const op = "mongodbsalesreturn.FindByID"

	var r entity.SalesReturn
	if err := d.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&r); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.SalesReturn{}, richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("sales_return.not_found")
		}
		return entity.SalesReturn{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return r, nil
}

// FindByOrderID lists every return already booked against one order, so a
// second return cannot give back more units than were sold.
func (d DB) FindByOrderID(ctx context.Context, orderID string) ([]entity.SalesReturn, error) {
	const op = "mongodbsalesreturn.FindByOrderID"

	cursor, err := d.collection.Find(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	returns := make([]entity.SalesReturn, 0)
	if err := cursor.All(ctx, &returns); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return returns, nil
}

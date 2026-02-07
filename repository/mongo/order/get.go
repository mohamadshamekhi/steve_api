package mongodborder

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) GetOrderByID(ctx context.Context, orderID string) (entity.Order, error) {
	const op = "mongodborder.GetOrderByID"

	filter := bson.M{"_id": orderID}

	var order entity.Order

	err := d.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Order{}, richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("order.not_found").
				WithMeta(map[string]any{"order_id": orderID})
		}
		return entity.Order{}, richerror.New(op).
			WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}

	return order, nil
}

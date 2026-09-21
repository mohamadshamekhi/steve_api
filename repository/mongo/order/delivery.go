package mongodborder

import (
	"context"
	"errors"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// EnsureDeliveryToken stores token on the order unless it already has one,
// and returns whichever token the order ends up with.
func (d DB) EnsureDeliveryToken(ctx context.Context, orderID string, token string) (string, error) {
	const op = "mongodborder.EnsureDeliveryToken"

	filter := bson.M{"_id": orderID, "delivery_token": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"delivery_token": token, "updated_at": time.Now()}}

	res, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	if res.ModifiedCount == 1 {
		return token, nil
	}

	order, gErr := d.GetOrderByID(ctx, orderID)
	if gErr != nil {
		return "", richerror.New(op).WithWarpError(gErr)
	}
	return order.DeliveryToken, nil
}

func (d DB) GetOrderByDeliveryToken(ctx context.Context, token string) (entity.Order, error) {
	const op = "mongodborder.GetOrderByDeliveryToken"

	var order entity.Order
	err := d.collection.FindOne(ctx, bson.M{"delivery_token": token}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Order{}, richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("order.delivery_not_found")
		}
		return entity.Order{}, richerror.New(op).
			WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}
	return order, nil
}

func (d DB) SetDeliveryConfirmedBy(ctx context.Context, orderID string, by string) error {
	const op = "mongodborder.SetDeliveryConfirmedBy"

	_, err := d.collection.UpdateOne(ctx, bson.M{"_id": orderID},
		bson.M{"$set": bson.M{"delivery_confirmed_by": by, "updated_at": time.Now()}})
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	return nil
}

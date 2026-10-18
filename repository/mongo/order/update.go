package mongodborder

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error {
	const op = "mongodborder.UpdateOrderStatus"
	filter := bson.M{"_id": orderID}

	now := time.Now()
	set := bson.M{
		"status":     status,
		"updated_at": now,
	}
	if status == entity.OrderStatusDelivered {
		set["delivered_at"] = now
	}
	update := bson.M{"$set": set}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("order.not_found").
			WithMeta(map[string]any{"order_id": orderID})
	}

	return nil
}

// SetCancelReason records why a cancelled order was cancelled. Called right
// after UpdateOrderStatus(..., OrderStatusCancelled).
func (d DB) SetCancelReason(ctx context.Context, orderID string, reason entity.CancelReason) error {
	const op = "mongodborder.SetCancelReason"
	filter := bson.M{"_id": orderID}
	update := bson.M{"$set": bson.M{"cancel_reason": reason, "updated_at": time.Now()}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("order.not_found").
			WithMeta(map[string]any{"order_id": orderID})
	}
	return nil
}

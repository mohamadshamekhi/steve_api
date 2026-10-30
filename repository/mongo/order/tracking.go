package mongodborder

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// SetTrackingCode records the courier tracking number of a shipped order.
// Called right before UpdateOrderStatus(..., OrderStatusShipped).
func (d DB) SetTrackingCode(ctx context.Context, orderID string, trackingCode string) error {
	const op = "mongodborder.SetTrackingCode"
	filter := bson.M{"_id": orderID}

	now := time.Now()
	update := bson.M{"$set": bson.M{
		"tracking_code": trackingCode,
		"shipped_at":    now,
		"updated_at":    now,
	}}

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

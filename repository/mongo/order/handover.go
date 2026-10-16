package mongodborder

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// SetPartnerHandover records the in-person handover of a partner order's goods
// and marks it delivered, in one update.
func (d DB) SetPartnerHandover(ctx context.Context, orderID, handoverBy string, handoverAt time.Time) error {
	const op = "mongodborder.SetPartnerHandover"
	filter := bson.M{"_id": orderID}
	update := bson.M{
		"$set": bson.M{
			"handover_by": handoverBy,
			"handover_at": handoverAt,
			"status":      entity.OrderStatusDelivered,
			"updated_at":  time.Now(),
		},
	}
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

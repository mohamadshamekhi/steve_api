package mongodborder

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// ApplyEdit writes a corrected order and appends its audit entry. It only
// matches while the order is still at prevUpdatedAt, so two admins editing
// the same order at once cannot silently overwrite each other.
func (d DB) ApplyEdit(ctx context.Context, order entity.Order, prevUpdatedAt time.Time, edit entity.OrderEdit) error {
	const op = "mongodborder.ApplyEdit"

	filter := bson.M{"_id": order.ID, "updated_at": prevUpdatedAt}
	update := bson.M{
		"$set": bson.M{
			"user_id":                order.UserID,
			"items":                  order.Items,
			"shipping_address":       order.ShippingAddress,
			"subtotal_amount":        order.SubtotalAmount,
			"shipping_amount":        order.ShippingAmount,
			"manual_discount_amount": order.ManualDiscountAmount,
			"final_amount":           order.FinalAmount,
			"updated_at":             order.UpdatedAt,
		},
		"$push": bson.M{"edits": edit},
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).WithKind(richerror.KindConflict).
			WithMessageKey("order.edit_conflict").
			WithMeta(map[string]any{"order_id": order.ID})
	}
	return nil
}

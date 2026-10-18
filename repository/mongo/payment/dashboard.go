package mongodbpayment

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// GetLatestStatusByOrderIDs returns, for each order id that has at least one
// payment, the status of its most recently created payment attempt. Orders
// with no payment record at all are simply absent from the map.
func (d DB) GetLatestStatusByOrderIDs(ctx context.Context, orderIDs []string) (map[string]entity.PaymentStatus, error) {
	const op = "mongodbpayment.GetLatestStatusByOrderIDs"
	if len(orderIDs) == 0 {
		return map[string]entity.PaymentStatus{}, nil
	}

	pipeline := []bson.M{
		{"$match": bson.M{"order_id": bson.M{"$in": orderIDs}}},
		{"$sort": bson.M{"created_at": -1}},
		{"$group": bson.M{"_id": "$order_id", "status": bson.M{"$first": "$status"}}},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var rows []struct {
		OrderID string               `bson:"_id"`
		Status  entity.PaymentStatus `bson:"status"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	out := make(map[string]entity.PaymentStatus, len(rows))
	for _, r := range rows {
		out[r.OrderID] = r.Status
	}
	return out, nil
}

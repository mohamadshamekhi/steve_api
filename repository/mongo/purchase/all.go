package mongodbpurchase

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PurchaseFilter narrows the purchase-invoice listing. Every field optional.
type PurchaseFilter struct {
	SupplierID   string
	Status       entity.PurchaseStatus
	PaymentState entity.PaymentState
	OrderID      string
	From         *time.Time
	To           *time.Time
}

func (d DB) All(ctx context.Context, f PurchaseFilter, p commonparam.PaginationRequest) ([]entity.Purchase, int64, error) {
	const op = "mongodbpurchase.All"

	filter := bson.M{}
	if f.SupplierID != "" {
		filter["supplier_id"] = f.SupplierID
	}
	if f.Status != "" {
		filter["status"] = f.Status
	}
	if f.PaymentState != "" {
		filter["payment_state"] = f.PaymentState
	}
	if f.OrderID != "" {
		filter["order_id"] = f.OrderID
	}
	if f.From != nil || f.To != nil {
		purchasedAt := bson.M{}
		if f.From != nil {
			purchasedAt["$gte"] = *f.From
		}
		if f.To != nil {
			purchasedAt["$lte"] = *f.To
		}
		filter["purchased_at"] = purchasedAt
	}

	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	findOptions := options.Find().
		SetSkip(int64(p.GetOffset())).
		SetLimit(int64(p.GetPageSize())).
		SetSort(bson.D{{Key: "purchased_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	purchases := make([]entity.Purchase, 0)
	if err := cursor.All(ctx, &purchases); err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return purchases, total, nil
}

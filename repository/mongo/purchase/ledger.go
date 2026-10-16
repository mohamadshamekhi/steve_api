package mongodbpurchase

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// SupplierLedger aggregates confirmed invoices per supplier: what was bought,
// what has been paid, and what is still owed.
//
// Draft and cancelled invoices are excluded — a draft is not a liability, and
// a cancelled one never was.
func (d DB) SupplierLedger(ctx context.Context, req accountingparam.SupplierLedgerRequest) ([]accountingparam.SupplierLedgerRow, error) {
	const op = "mongodbpurchase.SupplierLedger"

	match := bson.M{"status": entity.PurchaseStatusConfirmed}
	if req.SupplierID != "" {
		match["supplier_id"] = req.SupplierID
	}
	if req.From != nil || req.To != nil {
		purchasedAt := bson.M{}
		if req.From != nil {
			purchasedAt["$gte"] = *req.From
		}
		if req.To != nil {
			purchasedAt["$lte"] = *req.To
		}
		match["purchased_at"] = purchasedAt
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$group": bson.M{
			"_id":              "$supplier_id",
			"label":            bson.M{"$first": "$supplier_name"},
			"invoice_count":    bson.M{"$sum": 1},
			"total_purchased":  bson.M{"$sum": "$total_amount"},
			"total_paid":       bson.M{"$sum": "$paid_amount"},
			"last_purchase_at": bson.M{"$max": "$purchased_at"},
		}},
		{"$sort": bson.M{"total_purchased": -1}},
		{"$limit": 500},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var aggRows []struct {
		ID             string    `bson:"_id"`
		Label          string    `bson:"label"`
		InvoiceCount   int64     `bson:"invoice_count"`
		TotalPurchased uint64    `bson:"total_purchased"`
		TotalPaid      uint64    `bson:"total_paid"`
		LastPurchaseAt time.Time `bson:"last_purchase_at"`
	}
	if err := cursor.All(ctx, &aggRows); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	rows := make([]accountingparam.SupplierLedgerRow, 0, len(aggRows))
	for _, r := range aggRows {
		balance := int64(r.TotalPurchased) - int64(r.TotalPaid)
		if req.OnlyDebt && balance <= 0 {
			continue
		}
		row := accountingparam.SupplierLedgerRow{
			SupplierID:     r.ID,
			SupplierName:   r.Label,
			InvoiceCount:   r.InvoiceCount,
			TotalPurchased: r.TotalPurchased,
			TotalPaid:      r.TotalPaid,
			Balance:        balance,
		}
		if !r.LastPurchaseAt.IsZero() {
			last := r.LastPurchaseAt
			row.LastPurchaseAt = &last
		}
		rows = append(rows, row)
	}
	return rows, nil
}

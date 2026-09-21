package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// SalesReturn is a customer return (مرجوعی / برگشت از فروش) against a
// previously registered order. It is its own document rather than a mutation
// of the order: the order is what was sold, the return is what came back, and
// the profit report needs both.
type SalesReturn struct {
	ID      string            `bson:"_id" json:"id"`
	OrderID string            `bson:"order_id" json:"order_id"`
	UserID  string            `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Items   []SalesReturnItem `bson:"items" json:"items"`
	// RefundAmount is what the customer actually got back. It can be lower
	// than the items' value (کسر هزینه‌ی حمل یا جریمه).
	RefundAmount uint64        `bson:"refund_amount" json:"refund_amount"`
	ItemsAmount  uint64        `bson:"items_amount" json:"items_amount"`
	ReturnedCost uint64        `bson:"returned_cost" json:"returned_cost"`
	RefundMethod PaymentMethod `bson:"refund_method,omitempty" json:"refund_method,omitempty"`
	// IsRefunded is false when the money has not been handed back yet, so the
	// cashbox is not told about it.
	IsRefunded bool      `bson:"is_refunded" json:"is_refunded"`
	Reason     string    `bson:"reason,omitempty" json:"reason,omitempty"`
	Note       string    `bson:"note,omitempty" json:"note,omitempty"`
	ReturnedAt time.Time `bson:"returned_at" json:"returned_at"`
	CreatedBy  string    `bson:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

type SalesReturnItem struct {
	ProductID         string             `bson:"product_id,omitempty" json:"product_id,omitempty"`
	VariantID         string             `bson:"variant_id" json:"variant_id"`
	ProductName       string             `bson:"product_name,omitempty" json:"product_name,omitempty"`
	VariantAttributes []VariantAttribute `bson:"variant_attributes,omitempty" json:"variant_attributes,omitempty"`
	Quantity          uint64             `bson:"quantity" json:"quantity"`
	// UnitPrice/UnitCost are copied from the original order line so the
	// return is valued at what was actually sold, not at today's price.
	UnitPrice uint64 `bson:"unit_price" json:"unit_price"`
	UnitCost  uint64 `bson:"unit_cost" json:"unit_cost"`
	// Restock is false for goods that came back broken or that never lived in
	// the warehouse (a service, an externally supplied item).
	Restock bool `bson:"restock" json:"restock"`
}

func (r *SalesReturn) BeforeCreate() {
	r.ID = ulid.Make().String()
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	if r.ReturnedAt.IsZero() {
		r.ReturnedAt = now
	}
	r.Recalculate()
}

// Recalculate refreshes the value of the returned goods and their cost basis.
func (r *SalesReturn) Recalculate() {
	var items, cost uint64
	for _, item := range r.Items {
		items += item.Quantity * item.UnitPrice
		cost += item.Quantity * item.UnitCost
	}
	r.ItemsAmount = items
	r.ReturnedCost = cost
}

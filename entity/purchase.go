package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// PurchaseStatus is the lifecycle of a purchase invoice (فاکتور خرید).
type PurchaseStatus string

const (
	// PurchaseStatusDraft — پیش‌نویس؛ هنوز روی انبار اثر نگذاشته.
	PurchaseStatusDraft PurchaseStatus = "draft"
	// PurchaseStatusConfirmed — تأیید شده؛ کالا وارد انبار و هزینه ثبت شده.
	PurchaseStatusConfirmed PurchaseStatus = "confirmed"
	// PurchaseStatusCanceled — لغو شده؛ اگر تأیید شده بود، اثرش برگشت خورده.
	PurchaseStatusCanceled PurchaseStatus = "canceled"
)

// PaymentState tracks how much of the invoice the shop has actually paid.
type PaymentState string

const (
	PaymentStateUnpaid  PaymentState = "unpaid"
	PaymentStatePartial PaymentState = "partial"
	PaymentStatePaid    PaymentState = "paid"
)

// Purchase is a buy-side invoice: goods the shop bought from a supplier, to
// resell. Confirming it pushes every item into the inventory ledger and
// updates each variant's weighted-average cost.
type Purchase struct {
	ID            string         `bson:"_id" json:"id"`
	SupplierID    string         `bson:"supplier_id" json:"supplier_id"`
	SupplierName  string         `bson:"supplier_name,omitempty" json:"supplier_name,omitempty"`
	InvoiceNumber string         `bson:"invoice_number,omitempty" json:"invoice_number,omitempty"`
	Items         []PurchaseItem `bson:"items" json:"items"`
	TotalAmount   uint64         `bson:"total_amount" json:"total_amount"`
	PaidAmount    uint64         `bson:"paid_amount" json:"paid_amount"`
	PaymentState  PaymentState   `bson:"payment_state" json:"payment_state"`
	Status        PurchaseStatus `bson:"status" json:"status"`
	Note          string         `bson:"note,omitempty" json:"note,omitempty"`
	// OrderID links a purchase made specifically to fulfil a customer order
	// (the "bought it from next door for this customer" case), so the profit
	// report can pair the two.
	OrderID     string     `bson:"order_id,omitempty" json:"order_id,omitempty"`
	PurchasedAt time.Time  `bson:"purchased_at" json:"purchased_at"`
	ConfirmedAt *time.Time `bson:"confirmed_at,omitempty" json:"confirmed_at,omitempty"`
	CreatedBy   string     `bson:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`
}

type PurchaseItem struct {
	ProductID         string             `bson:"product_id" json:"product_id"`
	VariantID         string             `bson:"variant_id" json:"variant_id"`
	ProductName       string             `bson:"product_name,omitempty" json:"product_name,omitempty"`
	VariantAttributes []VariantAttribute `bson:"variant_attributes,omitempty" json:"variant_attributes,omitempty"`
	Quantity          uint64             `bson:"quantity" json:"quantity"`
	UnitCost          uint64             `bson:"unit_cost" json:"unit_cost"`
	TotalCost         uint64             `bson:"total_cost" json:"total_cost"`
	// AddToStock is false when the goods never enter the shop's inventory —
	// bought and handed straight to the customer. The cost is still recorded.
	AddToStock bool `bson:"add_to_stock" json:"add_to_stock"`
}

func (p *Purchase) BeforeCreate() {
	p.ID = ulid.Make().String()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.PurchasedAt.IsZero() {
		p.PurchasedAt = now
	}
	if p.Status == "" {
		p.Status = PurchaseStatusDraft
	}
	p.Recalculate()
}

// Recalculate refreshes per-item totals, the invoice total and the payment
// state. Call it after any change to Items or PaidAmount.
func (p *Purchase) Recalculate() {
	var total uint64
	for i := range p.Items {
		p.Items[i].TotalCost = p.Items[i].Quantity * p.Items[i].UnitCost
		total += p.Items[i].TotalCost
	}
	p.TotalAmount = total

	switch {
	case p.PaidAmount == 0:
		p.PaymentState = PaymentStateUnpaid
	case p.PaidAmount >= p.TotalAmount:
		p.PaymentState = PaymentStatePaid
	default:
		p.PaymentState = PaymentStatePartial
	}
}

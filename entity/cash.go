package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// PaymentMethod is how money actually moved: the till, a card reader, the
// online gateway, a bank transfer, or a credit arrangement.
type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "cash"     // نقد
	PaymentMethodCard     PaymentMethod = "card"     // کارتخوان / کارت به کارت
	PaymentMethodGateway  PaymentMethod = "gateway"  // درگاه پرداخت آنلاین
	PaymentMethodTransfer PaymentMethod = "transfer" // حواله / شبا
	PaymentMethodCredit   PaymentMethod = "credit"   // نسیه
	PaymentMethodOther    PaymentMethod = "other"    // سایر
)

func (m PaymentMethod) IsValid() bool {
	switch m {
	case PaymentMethodCash, PaymentMethodCard, PaymentMethodGateway,
		PaymentMethodTransfer, PaymentMethodCredit, PaymentMethodOther:
		return true
	}
	return false
}

// CashDirection says whether money came in or went out.
type CashDirection string

const (
	CashIn  CashDirection = "in"  // دریافت
	CashOut CashDirection = "out" // پرداخت
)

// Cash transaction reference kinds, mirroring the stock ledger's RefType.
const (
	CashRefOrder    = "order"
	CashRefPurchase = "purchase"
	CashRefExpense  = "expense"
	CashRefReturn   = "sales_return"
	CashRefManual   = "manual"
)

// CashTransaction is one line of the صندوق/دخل ledger: every rial that
// entered or left the shop, whatever the reason. It is append-only and is
// deliberately independent of the profit calculation — profit is accrual
// (ledger-based), the cashbox is cash-based, and the two legitimately differ
// (نسیه, پیش‌پرداخت, هزینه‌ی ثبت‌شده‌ی پرداخت‌نشده).
type CashTransaction struct {
	ID          string        `bson:"_id" json:"id"`
	Direction   CashDirection `bson:"direction" json:"direction"`
	Method      PaymentMethod `bson:"method" json:"method"`
	Amount      uint64        `bson:"amount" json:"amount"`
	Description string        `bson:"description,omitempty" json:"description,omitempty"`
	// RefType/RefID link the line back to what caused it, so an auto-recorded
	// line can be traced (and reversed) instead of looking like a mystery.
	RefType    string    `bson:"ref_type,omitempty" json:"ref_type,omitempty"`
	RefID      string    `bson:"ref_id,omitempty" json:"ref_id,omitempty"`
	SupplierID string    `bson:"supplier_id,omitempty" json:"supplier_id,omitempty"`
	IsAuto     bool      `bson:"is_auto,omitempty" json:"is_auto,omitempty"`
	OccurredAt time.Time `bson:"occurred_at" json:"occurred_at"`
	CreatedBy  string    `bson:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}

func (t *CashTransaction) BeforeCreate() {
	t.ID = ulid.Make().String()
	now := time.Now()
	t.CreatedAt = now
	if t.OccurredAt.IsZero() {
		t.OccurredAt = now
	}
	if t.Method == "" {
		t.Method = PaymentMethodCash
	}
}

// SignedAmount is the amount as it affects the balance: negative for payouts.
func (t CashTransaction) SignedAmount() int64 {
	if t.Direction == CashOut {
		return -int64(t.Amount)
	}
	return int64(t.Amount)
}

package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// StockMovementType is why stock moved. Every movement is append-only: the
// ledger is the source of truth for both انبارداری and حسابداری, while
// Variant.Quantity/AvgCost are a cache of the ledger's running totals.
type StockMovementType string

const (
	// StockMovementPurchaseIn — خرید از تأمین‌کننده و ورود به انبار.
	StockMovementPurchaseIn StockMovementType = "purchase_in"
	// StockMovementSaleOut — فروش از موجودی انبار.
	StockMovementSaleOut StockMovementType = "sale_out"
	// StockMovementReturnIn — برگشت به انبار (لغو یا مرجوعی سفارش).
	StockMovementReturnIn StockMovementType = "return_in"
	// StockMovementAdjustment — اصلاح دستی/انبارگردانی (+ یا -).
	StockMovementAdjustment StockMovementType = "adjustment"
	// StockMovementExternalSupply — فروش کالایی که در انبار نبود و در لحظه از
	// تأمین‌کننده تهیه شد. موجودی را تغییر نمی‌دهد ولی هزینه و سود را ثبت می‌کند.
	StockMovementExternalSupply StockMovementType = "external_supply"
	// StockMovementServiceSale — فروش خدمات؛ فقط برای حسابداری.
	StockMovementServiceSale StockMovementType = "service_sale"
)

func (t StockMovementType) IsValid() bool {
	switch t {
	case StockMovementPurchaseIn, StockMovementSaleOut, StockMovementReturnIn,
		StockMovementAdjustment, StockMovementExternalSupply, StockMovementServiceSale:
		return true
	}
	return false
}

// AffectsStock reports whether this movement type changes Variant.Quantity.
func (t StockMovementType) AffectsStock() bool {
	switch t {
	case StockMovementPurchaseIn, StockMovementSaleOut, StockMovementReturnIn, StockMovementAdjustment:
		return true
	}
	return false
}

// Movement reference kinds, stored in StockMovement.RefType.
const (
	StockRefOrder    = "order"
	StockRefPurchase = "purchase"
	StockRefManual   = "manual"
)

// StockMovement is one append-only line of the inventory/accounting ledger.
type StockMovement struct {
	ID                string             `bson:"_id" json:"id"`
	ProductID         string             `bson:"product_id" json:"product_id"`
	VariantID         string             `bson:"variant_id" json:"variant_id"`
	ProductName       string             `bson:"product_name,omitempty" json:"product_name,omitempty"`
	Kind              ProductKind        `bson:"kind,omitempty" json:"kind,omitempty"`
	Type              StockMovementType  `bson:"type" json:"type"`
	VariantAttributes []VariantAttribute `bson:"variant_attributes,omitempty" json:"variant_attributes,omitempty"`

	// Quantity is signed: positive for stock coming in, negative for stock
	// going out. For external_supply/service_sale it is negative too, to keep
	// "units sold" summable across every sale type.
	Quantity int64 `bson:"quantity" json:"quantity"`

	UnitCost   uint64 `bson:"unit_cost" json:"unit_cost"`                         // قیمت تمام‌شده هر واحد
	UnitPrice  uint64 `bson:"unit_price,omitempty" json:"unit_price,omitempty"`   // قیمت فروش هر واحد (برای فروش)
	TotalCost  uint64 `bson:"total_cost" json:"total_cost"`                       // |Quantity| × UnitCost
	TotalPrice uint64 `bson:"total_price,omitempty" json:"total_price,omitempty"` // |Quantity| × UnitPrice
	// Profit is TotalPrice - TotalCost for sale movements, 0 otherwise. Signed
	// because a sale below cost is a real, reportable loss.
	Profit int64 `bson:"profit" json:"profit"`

	RefType    string `bson:"ref_type,omitempty" json:"ref_type,omitempty"` // order | purchase | manual
	RefID      string `bson:"ref_id,omitempty" json:"ref_id,omitempty"`
	SupplierID string `bson:"supplier_id,omitempty" json:"supplier_id,omitempty"`

	// QuantityAfter/AvgCostAfter snapshot the variant's running totals right
	// after this movement, so a report never has to replay the whole ledger.
	QuantityAfter int64  `bson:"quantity_after" json:"quantity_after"`
	AvgCostAfter  uint64 `bson:"avg_cost_after" json:"avg_cost_after"`

	Note      string    `bson:"note,omitempty" json:"note,omitempty"`
	CreatedBy string    `bson:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

func (m *StockMovement) BeforeCreate() {
	m.ID = ulid.Make().String()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	abs := m.Quantity
	if abs < 0 {
		abs = -abs
	}
	m.TotalCost = uint64(abs) * m.UnitCost
	m.TotalPrice = uint64(abs) * m.UnitPrice
	if m.UnitPrice > 0 {
		m.Profit = int64(m.TotalPrice) - int64(m.TotalCost)
	}
}

// NewWeightedAverageCost folds an incoming purchase into the running weighted
// average. currentQty/currentAvg describe stock before the purchase.
//
// When stock is empty (or negative from backorders) the incoming price simply
// becomes the new average — there is nothing left to average against.
func NewWeightedAverageCost(currentQty int64, currentAvg uint64, inQty uint64, inCost uint64) uint64 {
	if inQty == 0 {
		return currentAvg
	}
	if currentQty <= 0 || currentAvg == 0 {
		return inCost
	}
	totalValue := uint64(currentQty)*currentAvg + inQty*inCost
	totalQty := uint64(currentQty) + inQty
	return totalValue / totalQty
}

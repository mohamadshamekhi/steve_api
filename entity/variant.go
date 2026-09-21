package entity

import (
	"github.com/oklog/ulid/v2"
)

type VariantAttribute struct {
	Key   string `bson:"key" json:"key"`               // ex: "color"
	Value string `bson:"value" json:"value"`           // ex: "red"
	Label string `bson:"label,omitempty" json:"label"` // ex:"قزمز",
	Meta  string `bson:"meta,omitempty" json:"meta"`
}

// StockMode says whether this variant is backed by real shop inventory or is
// sourced on the spot from another shop when a customer orders it.
type StockMode string

const (
	// StockModeTracked — کالایی که در مغازه دارم؛ موجودی شمرده و کم می‌شود.
	StockModeTracked StockMode = "tracked"
	// StockModeOnDemand — کالایی که در مغازه ندارم؛ سفارش که آمد از تأمین‌کننده
	// می‌خرم. موجودی شمرده نمی‌شود و همیشه قابل فروش است.
	StockModeOnDemand StockMode = "on_demand"
	// StockModeNone — خدمات؛ نه موجودی دارد نه تأمین.
	StockModeNone StockMode = "none"
)

func (m StockMode) IsValid() bool {
	return m == StockModeTracked || m == StockModeOnDemand || m == StockModeNone
}

type Variant struct {
	ID               string             `bson:"id" json:"id"`
	Price            uint64             `bson:"price" json:"price"`
	Quantity         uint64             `bson:"quantity" json:"quantity"`
	VariantAttribute []VariantAttribute `bson:"variant_attribute,omitempty" json:"variant_attribute,omitempty"`
	Warranty         string             `bson:"warranty" json:"warranty"`

	// --- inventory & accounting ---

	// StockMode. Empty is read as StockModeTracked so documents written before
	// the inventory feature keep their old behaviour.
	StockMode StockMode `bson:"stock_mode,omitempty" json:"stock_mode,omitempty"`
	// AllowBackorder lets the variant keep selling after Quantity hits zero by
	// sourcing the extra units externally — the "I have 10 in the shop, the
	// 11th I buy from next door" case. Only meaningful for StockModeTracked.
	AllowBackorder bool `bson:"allow_backorder,omitempty" json:"allow_backorder,omitempty"`
	// BackorderLimit caps how many units past Quantity may be sold. nil means
	// unlimited (when AllowBackorder is true).
	BackorderLimit *uint64 `bson:"backorder_limit,omitempty" json:"backorder_limit,omitempty"`
	// AvgCost is the weighted-average purchase price of one unit currently in
	// stock. Maintained by the inventory ledger on every purchase_in.
	AvgCost uint64 `bson:"avg_cost,omitempty" json:"avg_cost,omitempty"`
	// SupplierID is the default supplier this variant is sourced from when it
	// has to be supplied externally (on-demand or backorder).
	SupplierID string `bson:"supplier_id,omitempty" json:"supplier_id,omitempty"`
	// SupplyCost is the expected unit cost when sourcing externally. Used as
	// the fallback cost basis for profit reporting until the real purchase
	// invoice is registered.
	SupplyCost uint64 `bson:"supply_cost,omitempty" json:"supply_cost,omitempty"`
}

func (v *Variant) BeforeCreate() {
	v.ID = ulid.Make().String()
	if v.StockMode == "" {
		v.StockMode = StockModeTracked
	}
}

// GetStockMode normalises the zero value of pre-inventory documents.
func (v Variant) GetStockMode() StockMode {
	if v.StockMode == "" {
		return StockModeTracked
	}
	return v.StockMode
}

// IsStockTracked reports whether selling this variant should decrement stock.
func (v Variant) IsStockTracked() bool {
	return v.GetStockMode() == StockModeTracked
}

// CanSell reports whether qty units may be sold right now, either out of
// stock or by sourcing them externally.
func (v Variant) CanSell(qty uint64) bool {
	switch v.GetStockMode() {
	case StockModeOnDemand, StockModeNone:
		return true
	default:
		if qty <= v.Quantity {
			return true
		}
		if !v.AllowBackorder {
			return false
		}
		if v.BackorderLimit == nil {
			return true
		}
		return qty-v.Quantity <= *v.BackorderLimit
	}
}

// SplitSupply splits qty into the part served from shop stock and the part
// that has to be sourced externally.
func (v Variant) SplitSupply(qty uint64) (fromStock uint64, external uint64) {
	if !v.IsStockTracked() {
		return 0, qty
	}
	if qty <= v.Quantity {
		return qty, 0
	}
	return v.Quantity, qty - v.Quantity
}

// UnitCost is the cost basis used for profit reporting: the weighted-average
// stock cost when there is one, otherwise the expected external supply cost.
func (v Variant) UnitCost() uint64 {
	if v.AvgCost > 0 {
		return v.AvgCost
	}
	return v.SupplyCost
}

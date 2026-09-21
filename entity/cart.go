package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type CartStatus string

const (
	CartActive     CartStatus = "active"
	CartInactive   CartStatus = "inactive"
	CartPaid       CartStatus = "paid"
	CartCanceled   CartStatus = "cancelled"
	PendingPayment CartStatus = "pending_payment"
)

type Cart struct {
	ID string `json:"id" bson:"_id"`

	// ManualDiscount is a discount given on the invoice as a whole at the
	// counter, in Toman. It is applied after per-line discounts, to whatever
	// is left of the total.
	ManualDiscount uint `bson:"manual_discount,omitempty"`

	UserID string `bson:"user_id,omitempty"`

	SessionID string `bson:"session_id,omitempty"`

	Items []CartItem `bson:"items"`

	Status CartStatus `bson:"status"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type CartItem struct {
	ProductID string `bson:"product_id"`

	VariantID string `bson:"variant_id"`

	Quantity int `bson:"quantity"`

	ProductName string `bson:"product_name"`

	VariantAttributes []VariantAttribute `bson:"variant_attributes"`

	PriceAtTime int64 `bson:"price_at_time"`

	// --- counter (in-person) invoice lines ---

	// IsCustom marks a line typed by hand at the counter — a repair fee, an
	// Apple ID, a part bought from the shop next door — that has no catalog
	// product behind it. Such a line has an empty ProductID and a generated
	// VariantID that only identifies it inside this cart.
	IsCustom bool `bson:"is_custom,omitempty"`

	// Kind is carried only by custom lines; a catalog line resolves its kind
	// from the variant at checkout time.
	Kind ProductKind `bson:"kind,omitempty"`

	// UnitCost is what one unit of a hand-added part cost the shop, so the
	// profit report stays honest for goods that never entered the warehouse.
	UnitCost uint `bson:"unit_cost,omitempty"`

	// ManualDiscount is the discount the salesperson granted on this line at
	// the counter, in Toman, for the whole line (not per unit). It is applied
	// on top of any campaign discount the line already earned.
	ManualDiscount uint `bson:"manual_discount,omitempty"`
}

// GetKind normalises a cart line whose kind was never set.
func (i CartItem) GetKind() ProductKind {
	if i.Kind == "" {
		return ProductKindPhysical
	}
	return i.Kind
}

// LineTotal is the gross price of the line, before any discount.
func (i CartItem) LineTotal() int64 {
	return i.PriceAtTime * int64(i.Quantity)
}

func (c *Cart) BeforeCreate() {
	c.ID = ulid.Make().String()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.Status = CartActive
	if c.Items == nil {
		c.Items = make([]CartItem, 0)
	}
}

func (c *Cart) CalculateSubtotal() uint {
	var total uint = 0
	for _, item := range c.Items {
		total += uint(item.PriceAtTime) * uint(item.Quantity)
	}
	return total
}

type CheckoutItemResult struct {
	CartItem        `json:",inline"`
	HasStockIssue   bool  `json:"has_stock_issue"`
	HasPriceChanged bool  `json:"has_price_changed"`
	CurrentPrice    int64 `json:"current_price"`
	DiscountAmount  uint  `json:"discount_amount"`

	// --- sourcing & cost, resolved while validating the cart ---

	// Kind marks a service line, which holds no stock and ships nothing.
	Kind ProductKind `json:"kind,omitempty"`
	// CostAtTime is the unit cost basis frozen for this checkout.
	CostAtTime uint `json:"cost_at_time,omitempty"`
	// Supply says how the units will be sourced.
	Supply SupplySource `json:"supply,omitempty"`
	// StockQuantity / ExternalQuantity split Quantity across shop stock and
	// an external supplier (the "11th unit comes from next door" case).
	StockQuantity    uint   `json:"stock_quantity,omitempty"`
	ExternalQuantity uint   `json:"external_quantity,omitempty"`
	SupplierID       string `json:"supplier_id,omitempty"`
}

type CheckoutSummary struct {
	Items              []CheckoutItemResult `bson:"items" json:"items,omitempty"`
	Subtotal           int64                `json:"subtotal,omitempty"`
	AppliedDiscounts   []AppliedDiscount    `json:"applied_discounts,omitempty"`
	FinalPrice         int64                `json:"final_price,omitempty"`
	IsReady            bool                 `json:"is_ready,omitempty"`
	DiscountedSubtotal uint                 `json:"discounted_subtotal,omitempty"`
}

func CanTransitionCartStatus(from, to CartStatus) bool {
	allowed := map[CartStatus][]CartStatus{
		CartActive: {
			CartInactive,
			CartPaid,
			CartCanceled,
			PendingPayment,
		},

		CartInactive: {
			CartActive,
			CartCanceled,
		},

		PendingPayment: {
			CartPaid,
			CartCanceled,
			CartActive,
		},

		CartPaid: {},

		CartCanceled: {},
	}

	for _, st := range allowed[from] {
		if st == to {
			return true
		}
	}

	return false
}

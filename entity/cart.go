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

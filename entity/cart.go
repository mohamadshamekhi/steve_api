package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type CartStatus string

const (
	CartActive   CartStatus = "active"
	CartInactive CartStatus = "inactive"
	CartPaid     CartStatus = "paid"
	CartCanceled CartStatus = "cancelled"
)

type Cart struct {
	ID string `json:"id"`

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
	c.Status = "active"
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
	CartItem
	HasStockIssue   bool
	HasPriceChanged bool
	CurrentPrice    int64
	DiscountAmount  uint
}

type CheckoutSummary struct {
	Items              []CheckoutItemResult
	Subtotal           int64
	AppliedDiscounts   []AppliedDiscount
	FinalPrice         int64
	IsReady            bool
	DiscountedSubtotal uint
}

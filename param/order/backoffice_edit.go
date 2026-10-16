package orderparam

import "github.com/abolfazlnorzad/steve/entity"

// EditOrderRequest corrects a registered order. Every field is optional;
// only what is sent changes. Items, when sent, is the complete new list of
// lines (a line left out is removed).
type EditOrderRequest struct {
	OrderID  string `param:"id" json:"-"`
	EditedBy string `json:"-"`

	// Reason is required: why the order is being corrected.
	Reason string `json:"reason"`

	UserID               *string         `json:"user_id"`
	ShippingAddress      *entity.Address `json:"shipping_address"`
	ShippingAmount       *uint           `json:"shipping_amount"`
	ManualDiscountAmount *uint           `json:"manual_discount_amount"`
	Items                []EditOrderItem `json:"items"`

	// RecordCashDifference books the change in final amount in the cashbox
	// (received or paid back). Defaults to true for a completed counter sale
	// (in_person, delivered) and false otherwise.
	RecordCashDifference *bool `json:"record_cash_difference"`
	// PaymentMethod for that cashbox line; defaults to cash.
	PaymentMethod entity.PaymentMethod `json:"payment_method"`
}

// EditOrderItem is one line of the corrected order.
//
// A catalog line is identified by variant_id (+ product_id for a new one).
// A hand-typed line sets is_custom with product_name/kind/unit_cost; to keep
// an existing custom line send its variant_id.
type EditOrderItem struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
	Quantity  uint   `json:"quantity"`
	// UnitPrice defaults to the line's current price (or the variant's price
	// for a new catalog line).
	UnitPrice *uint `json:"unit_price"`
	// DiscountAmount is the whole discount on the line, in Toman. Left out,
	// an existing line keeps its discount (scaled to the new quantity) and a
	// new line gets none.
	DiscountAmount *uint `json:"discount_amount"`

	IsCustom    bool               `json:"is_custom"`
	ProductName string             `json:"product_name"`
	Kind        entity.ProductKind `json:"kind"`
	UnitCost    *uint              `json:"unit_cost"`
}

type EditOrderResponse struct {
	Order entity.Order     `json:"order"`
	Edit  entity.OrderEdit `json:"edit"`
}

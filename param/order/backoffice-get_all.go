package orderparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllOrderRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     BackofficeFilterOrder         `json:"filter" ,query:"filter"`
}

type BackofficeFilterOrder struct {
	UserID *string `json:"user_id" query:"user_id"`
	// Status is one status or a comma-separated list (processing,shipped).
	Status *string `json:"status" query:"status"`
	// OrderID matches an exact order ID.
	OrderID *string `json:"order_id" query:"order_id"`
	// Phone matches any part of the customer's phone number.
	Phone          *string `json:"phone" query:"phone"`
	PurchaseType   *string `json:"purchase_type" query:"purchase_type"`
	IsPartnerOrder *bool   `json:"is_partner_order" query:"is_partner_order"`
	// From / To are YYYY-MM-DD in Tehran time, both inclusive, on created_at.
	From      *string `json:"from" query:"from"`
	To        *string `json:"to" query:"to"`
	MinAmount *uint   `json:"min_amount" query:"min_amount"`
	MaxAmount *uint   `json:"max_amount" query:"max_amount"`

	// Resolved by the handler from From / To.
	FromTime *time.Time `json:"-" query:"-"`
	ToTime   *time.Time `json:"-" query:"-"`
}

type BackofficeOrderList struct {
	ID              string             `json:"id" bson:"_id"`
	UserPhoneNumber string             `json:"user_phone_number" bson:"user_phone_number"`
	UserName        string             `json:"user_name" bson:"user_name"`
	UserID          string             `json:"user_id"`
	Items           []entity.OrderItem `json:"items"`
	ShippingAddress entity.Address     `json:"shipping_address"`
	SubtotalAmount  uint               `json:"subtotal_amount"`
	ShippingAmount  uint               `json:"shipping_amount"`
	FinalAmount     uint               `json:"final_amount"`
	// AppliedDiscounts are the campaign/coupon discounts on the order.
	AppliedDiscounts []entity.AppliedDiscount `json:"applied_discounts"`
	// TotalDiscountAmount is everything taken off the invoice:
	// subtotal + shipping - final (coupon, campaign and manual discounts).
	TotalDiscountAmount uint `json:"total_discount_amount"`
	// ManualDiscount is the whole-invoice discount given at the counter.
	ManualDiscount uint            `json:"manual_discount"`
	Shipping       entity.Shipping `json:"shipping,omitempty"`
	// IsPostPayShipping: cart had a freight-only (باربری) item, so shipping
	// fee is collected on delivery and ShippingAmount above is 0.
	IsPostPayShipping bool               `json:"is_post_pay_shipping,omitempty"`
	Payments          []entity.Payment   `json:"payments,omitempty"`
	Status            entity.OrderStatus `json:"status"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	// IsEdited / Edits: super-admin corrections of this order, oldest first.
	IsEdited bool               `json:"is_edited"`
	Edits    []entity.OrderEdit `json:"edits,omitempty"`
}

type GetAllOrderResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Orders     []BackofficeOrderList          `json:"orders" bson:"orders"`
}

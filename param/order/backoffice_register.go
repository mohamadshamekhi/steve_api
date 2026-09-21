package orderparam

import "github.com/abolfazlnorzad/steve/entity"

type RegisterOrderBackofficeRequest struct {
	UserID       string `json:"user_id"`
	DiscountCode string `json:"discount_code"`
	// ManualDiscountAmount is the discount given by hand on the invoice as a
	// whole, in Toman. Leave it out to use whatever was already set on the
	// cart; send it to override that value for this registration.
	ManualDiscountAmount *uint `json:"manual_discount_amount"`
}

type RegisterOrderBackofficeResponse struct {
	Order BackofficeOrderInfo `json:"order"`
}

type BackofficeOrderInfo struct {
	ID                  string                   `json:"id"`
	UserID              string                   `json:"user_id"`
	CartID              string                   `json:"cart_id"`
	Items               []OrderItemInfo          `json:"items"`
	AppliedDiscounts    []entity.AppliedDiscount `json:"applied_discounts,omitempty"`
	SubtotalAmount      uint                     `json:"subtotal_amount"`
	TotalDiscountAmount uint                     `json:"total_discount_amount,omitempty"`
	// ItemsManualDiscount is the part of TotalDiscountAmount given by hand on
	// individual lines; ManualDiscount is the part given on the whole invoice.
	ItemsManualDiscount uint               `json:"items_manual_discount,omitempty"`
	ManualDiscount      uint               `json:"manual_discount,omitempty"`
	FinalAmount         uint               `json:"final_amount"`
	Status              entity.OrderStatus `json:"status"`
}

type OrderItemInfo struct {
	ProductID         string                    `json:"product_id"`
	VariantID         string                    `json:"variant_id"`
	Quantity          uint                      `json:"quantity"`
	ProductName       string                    `json:"product_name,omitempty"`
	VariantAttributes []entity.VariantAttribute `json:"variant_attributes"`
	TotalPriceAtTime  uint                      `json:"total_price_at_time"`
	DiscountAmount    uint                      `json:"discount_amount,omitempty"`
	ManualDiscount    uint                      `json:"manual_discount,omitempty"`
	IsCustom          bool                      `json:"is_custom,omitempty"`
	Kind              entity.ProductKind        `json:"kind,omitempty"`
}

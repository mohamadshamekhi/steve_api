package orderparam

import "github.com/abolfazlnorzad/steve/entity"

type RegisterOrderBackofficeRequest struct {
	UserID       string `json:"user_id"`
	DiscountCode string `json:"discount_code"`
}

type RegisterOrderBackofficeResponse struct {
	Order BackofficeOrderInfo `json:"order"`
}

type BackofficeOrderInfo struct {
	ID                  string                   `json:"id"`
	UserID              string                   `json:"user_id"`
	CartID              string                   `json:"cart_id"`
	Items               []OrderItemInfo          `json:"items"`
	AppliedDiscounts    []entity.AppliedDiscount `json:"applied_discounts"`
	SubtotalAmount      uint                     `json:"subtotal_amount"`
	TotalDiscountAmount uint                     `json:"total_discount_amount"`
	FinalAmount         uint                     `json:"final_amount"`
	Status              entity.OrderStatus       `json:"status"`
}

type OrderItemInfo struct {
	ProductID         string                    `json:"product_id"`
	VariantID         string                    `json:"variant_id"`
	Quantity          uint                      `json:"quantity"`
	VariantAttributes []entity.VariantAttribute `json:"variant_attributes"`
	TotalPriceAtTime  uint                      `json:"total_price_at_time"`
	DiscountAmount    uint                      `json:"discount_amount"`
}

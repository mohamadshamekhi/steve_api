package cartparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
)

type MyActiveCartBackofficeRequest struct {
	UserID string `json:"user_id"`
}

type MyActiveCartBackofficeResponse struct {
	Cart BackofficeCartDetail `json:"cart"`
}

type BackofficeCartDetail struct {
	ID        string                     `json:"id"`
	Status    string                     `json:"status"`
	UserID    string                     `json:"user_id"`
	Items     []BackofficeCartItemDetail `json:"items"`
	CreatedAt time.Time                  `json:"created_at"`
	UpdatedAt time.Time                  `json:"updated_at"`
	Subtotal  uint                       `json:"subtotal"`
	// ManualDiscount is the counter discount on the invoice as a whole.
	ManualDiscount uint `json:"manual_discount"`
	// ItemsManualDiscount is the net sum of per-line adjustments (positive =
	// discount, negative = markup).
	ItemsManualDiscount int64 `json:"items_manual_discount"`
	// PayableAmount is Subtotal minus line adjustments and invoice discount,
	// i.e. what the invoice shows before campaign discounts at registration.
	PayableAmount int64 `json:"payable_amount"`
}
type BackofficeCartItemDetail struct {
	ProductID         string                    `json:"product_id"`
	VariantID         string                    `json:"variant_id"`
	Quantity          int                       `json:"quantity"`
	ProductName       string                    `json:"product_name"`
	VariantAttributes []entity.VariantAttribute `json:"variant_attributes"`
	PriceAtTime       int64                     `json:"price_at_time"`
	Url               string                    `json:"url"`
	IsCustom          bool                      `json:"is_custom,omitempty"`
	Kind              entity.ProductKind        `json:"kind,omitempty"`
	UnitCost          uint                      `json:"unit_cost,omitempty"`
	ManualDiscount    int64                     `json:"manual_discount,omitempty"`
}

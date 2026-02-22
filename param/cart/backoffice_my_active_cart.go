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
}
type BackofficeCartItemDetail struct {
	ProductID         string                    `json:"product_id"`
	VariantID         string                    `json:"variant_id"`
	Quantity          int                       `json:"quantity"`
	ProductName       string                    `json:"product_name"`
	VariantAttributes []entity.VariantAttribute `json:"variant_attributes"`
	PriceAtTime       int64                     `json:"price_at_time"`
	Url               string                    `json:"url"`
}

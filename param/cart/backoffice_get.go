package cartparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
)

type GetCartRequest struct {
	ID string `json:"id"`
}

type GetCartResponse struct {
	Cart CartDetail `json:"cart"`
}

type CartDetail struct {
	ID string `json:"id"`

	Status string `json:"status"`

	User struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
	} `json:"user"`

	Items []CartItemDetail `json:"items"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CartItemDetail struct {
	ProductID         string                    `json:"product_id"`
	VariantID         string                    `json:"variant_id"`
	Quantity          int                       `bson:"quantity"`
	ProductName       string                    `bson:"product_name"`
	VariantAttributes []entity.VariantAttribute `bson:"variant_attributes"`
	PriceAtTime       int64                     `bson:"price_at_time"`
	Url               string                    `json:"url"`
	IsCustom          bool                      `json:"is_custom,omitempty"`
	Kind              entity.ProductKind        `json:"kind,omitempty"`
	ManualDiscount    uint                      `json:"manual_discount,omitempty"`
}

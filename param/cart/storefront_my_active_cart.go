package cartparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
)

type MyActiveCartRequest struct {
	UserID string `json:"user_id"`
}

type MyActiveCartResponse struct {
	Cart StoreFrontCartDetail `json:"cart"`
}

type StoreFrontCartDetail struct {
	ID string `json:"id"`

	Status string `json:"status"`

	User struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
	} `json:"user"`

	Items []StoreFrontCartItemDetail `json:"items"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Subtotal  uint      `json:"subtotal"`
}
type StoreFrontCartItemDetail struct {
	ProductID         string                    `json:"product_id"`
	VariantID         string                    `json:"variant_id"`
	Quantity          int                       `bson:"quantity"`
	ProductName       string                    `bson:"product_name"`
	VariantAttributes []entity.VariantAttribute `bson:"variant_attributes"`
	PriceAtTime       int64                     `bson:"price_at_time"`
	Url               string                    `json:"url"`
}

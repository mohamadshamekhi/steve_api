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
	Status *string `json:"status" query:"status"`
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
	Shipping        entity.Shipping    `json:"shipping,omitempty"`
	Status          entity.OrderStatus `json:"status"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type GetAllOrderResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Orders     []BackofficeOrderList          `json:"orders" bson:"orders"`
}

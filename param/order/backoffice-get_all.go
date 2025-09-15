package orderparam

import (
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
	Status          entity.OrderStatus `json:"status" bson:"status"`
	UserID          string             `json:"user_id" bson:"user_id"`
	UserPhoneNumber string             `json:"user_phone_number" bson:"user_phone_number"`
	UserName        string             `json:"user_name" bson:"user_name"`
	SubtotalAmount  uint               `bson:"subtotal_amount"`
	FinalAmount     uint               `bson:"final_amount"`
}

type GetAllOrderResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Orders     []BackofficeOrderList          `json:"orders" bson:"orders"`
}

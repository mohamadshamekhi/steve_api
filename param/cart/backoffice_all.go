package cartparam

import (
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllCartsRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     BackofficeFilterCart          `json:"filter" ,query:"filter"`
}

type BackofficeFilterCart struct {
	UserID *string `json:"user_id" query:"user_id"`
	Status *string `json:"status" query:"status"`
}

type BackofficeCartList struct {
	ID              string            `json:"id" bson:"_id"`
	Status          entity.CartStatus `json:"status" bson:"status"`
	UserID          string            `json:"user_id" bson:"user_id"`
	UserPhoneNumber string            `json:"user_phone_number" bson:"user_phone_number"`
	UserName        string            `json:"user_name" bson:"user_name"`
}

type GetAllCartsResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Carts      []BackofficeCartList
}

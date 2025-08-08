package userparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllUsersRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}

type GetAllUsersResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Users      []UserInfo
}

type UserInfo struct {
	ID          string           `json:"_id"`
	Name        string           `json:"name"`
	PhoneNumber string           `json:"phone_number"`
	CreatedAt   time.Time        `json:"created_at"`
	Addresses   []entity.Address `json:"addresses,omitempty"`
}

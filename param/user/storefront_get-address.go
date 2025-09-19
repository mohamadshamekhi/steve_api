package userparam

import "github.com/abolfazlnorzad/steve/entity"

type FindAddressByIDRequest struct {
	AddressID string `json:"address_id"`
	UserID    string `json:"user_id"`
}
type FindAddressByIDResponse struct {
	Address entity.Address `json:"address"`
}

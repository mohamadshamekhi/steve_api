package sellparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllMySellParamRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllMySellParamResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Sells      []SellInfo                     `json:"sells"`
}

type SellInfo struct {
	ID   string `json:"id" bson:"_id"`
	User struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
	} `json:"user_id"`
	BrandName         string            `json:"brand_name"`
	ProductName       string            `json:"product_name"`
	Color             string            `json:"color"`
	BatteryHealth     string            `json:"battery_health"`
	RegisterStatus    string            `json:"register_status"`
	PhysicalCondition string            `json:"physical_condition"`
	Description       string            `json:"description"`
	Status            entity.SellStatus `json:"status"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type GetUserSellStaticsRequest struct {
	UserID string `json:"user_id"`
}

type GetUserSellStaticsResponse struct {
	Count uint64 `json:"count"`
}

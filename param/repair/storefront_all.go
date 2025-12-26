package repairparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllMyRepairParamRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllMyRepairParamResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Repairs    []RepairInfo                   `json:"repairs"`
}

type RepairInfo struct {
	ID   string `json:"id" bson:"_id"`
	User struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
	} `json:"user_id"`
	BrandName   string                 `json:"brand_name"`
	ProductName string                 `json:"product_name"`
	Description string                 `json:"description"`
	Status      entity.RepairStatus    `json:"status"`
	Galleries   []entity.RepairGallery `json:"galleries"`
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at" json:"updated_at"`
}

type GetUserReairStaticRequest struct {
	UserID string `json:"user_id"`
}

type GetUserReairStaticResponse struct {
	Count uint64 `json:"count"`
}

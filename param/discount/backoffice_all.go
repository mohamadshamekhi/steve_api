package discountparam

import (
	"time"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetAllDiscountsParamRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
	Filter     commonparam.FilterRequest     `json:"filter"`
}
type GetAllDiscountsParamResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Discounts  []DiscountsList                `json:"discounts"`
}

type DiscountsList struct {
	ID              string    `json:"id" bson:"_id"`
	Code            string    `json:"code"`
	Type            string    `json:"type"`
	Value           uint      `json:"value"`
	StartsAt        time.Time `json:"starts_at" bson:"starts_at"`
	ExpiresAt       time.Time `json:"expires_at" bson:"expires_at"`
	TotalUsageLimit int       `json:"total_usage_limit,omitempty"`
	UsageCount      int       `json:"usage_count"`
}

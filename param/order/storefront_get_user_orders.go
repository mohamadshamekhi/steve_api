package orderparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetUserOrderRequest struct {
	UserID     string                        `json:"user_id"`
	Pagination commonparam.PaginationRequest `json:"pagination" ,query:"pagination"`
}
type StoreFrontOrderInfo struct {
	ID              string             `json:"_id"`
	UserID          string             `json:"user_id"`
	Items           []entity.OrderItem `json:"items"`
	ShippingAddress entity.Address     `json:"shipping_address"`
	SubtotalAmount  uint               `json:"subtotal_amount"`
	ShippingAmount  uint               `json:"shipping_amount"`
	FinalAmount     uint               `json:"final_amount"`
	Shipping        entity.Shipping    `json:"shipping,omitempty"`
	Status          entity.OrderStatus `json:"status"`
	TrackingCode    string             `json:"tracking_code,omitempty"`
	ShippedAt       *time.Time         `json:"shipped_at,omitempty"`
	Payments        []entity.Payment   `json:"payments"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type GetUserOrderResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Orders     []StoreFrontOrderInfo          `json:"orders"`
}

type GetUserOrderStaticsRequest struct {
	UserID string `json:"user_id"`
}

type GetUserOrderStaticsResponse struct {
	OrderProcessing uint64 `json:"order_processing"`
	OrderDelivered  uint64 `json:"order_delivered"`
	OrderRefunded   uint64 `json:"order_refunded"`
}

type GetOrderByIDRequest struct {
	ID string `json:"id"`
}
type GetOrderByIDResponse struct {
	Order entity.Order `json:"order"`
}

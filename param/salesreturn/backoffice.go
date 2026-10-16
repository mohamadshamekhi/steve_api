package salesreturnparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type CreateSalesReturnItem struct {
	VariantID string `json:"variant_id"`
	Quantity  uint64 `json:"quantity"`
	// Restock defaults to true for stock-tracked goods; send false for items
	// that came back unsellable.
	Restock *bool `json:"restock"`
}

type CreateSalesReturnRequest struct {
	OrderID string                  `json:"order_id"`
	Items   []CreateSalesReturnItem `json:"items"`
	// RefundAmount overrides the items' value when the shop keeps part of the
	// money. Zero means "refund exactly what the items were sold for".
	RefundAmount *uint64              `json:"refund_amount"`
	RefundMethod entity.PaymentMethod `json:"refund_method"`
	IsRefunded   *bool                `json:"is_refunded"`
	Reason       string               `json:"reason"`
	Note         string               `json:"note"`
	ReturnedAt   *time.Time           `json:"returned_at"`
}

type CreateSalesReturnResponse struct {
	SalesReturn entity.SalesReturn `json:"sales_return"`
}

type SalesReturnFilter struct {
	OrderID string
	UserID  string
	From    *time.Time
	To      *time.Time
}

type ListSalesReturnsRequest struct {
	Filter     SalesReturnFilter
	Pagination commonparam.PaginationRequest
}

type ListSalesReturnsResponse struct {
	SalesReturns []entity.SalesReturn           `json:"sales_returns"`
	Summary      SalesReturnSummary             `json:"summary"`
	Pagination   commonparam.PaginationResponse `json:"pagination"`
}

type GetSalesReturnResponse struct {
	SalesReturn entity.SalesReturn `json:"sales_return"`
}

// SalesReturnSummary feeds both the list header and the profit-and-loss
// statement's «برگشت از فروش» line.
type SalesReturnSummary struct {
	Count        int64  `json:"count"`
	RefundAmount uint64 `json:"refund_amount"`
	ItemsAmount  uint64 `json:"items_amount"`
	ReturnedCost uint64 `json:"returned_cost"`
}

type MarkRefundedRequest struct {
	ID           string               `json:"-"`
	RefundMethod entity.PaymentMethod `json:"refund_method"`
}

type MarkRefundedResponse struct {
	SalesReturn entity.SalesReturn `json:"sales_return"`
}

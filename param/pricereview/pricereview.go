package pricereviewparam

import (
	"time"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type GetSettingRequest struct{}
type GetSettingResponse struct {
	IntervalDays int `json:"interval_days"`
}

type UpdateSettingRequest struct {
	IntervalDays int `json:"interval_days"`
}
type UpdateSettingResponse struct{}

type ListNeedsReviewRequest struct {
	Pagination commonparam.PaginationRequest
}
type ListNeedsReviewResponse struct {
	Items []NeedsReviewItem `json:"items"`
	Total int64             `json:"total"`
}

type NeedsReviewItem struct {
	ProductID      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	CurrentPrice   uint64    `json:"current_price"`
	LastReviewedAt time.Time `json:"last_reviewed_at"`
	DaysOverdue    int       `json:"days_overdue"`
}

type ConfirmRequest struct {
	ProductID string `json:"product_id"`
}
type ConfirmResponse struct{}

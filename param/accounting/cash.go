package accountingparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type CreateCashTransactionRequest struct {
	Direction   entity.CashDirection `json:"direction"`
	Method      entity.PaymentMethod `json:"method"`
	Amount      uint64               `json:"amount"`
	Description string               `json:"description"`
	SupplierID  string               `json:"supplier_id"`
	RefType     string               `json:"ref_type"`
	RefID       string               `json:"ref_id"`
	OccurredAt  *time.Time           `json:"occurred_at"`
}

type CreateCashTransactionResponse struct {
	Transaction entity.CashTransaction `json:"transaction"`
}

type CashFilter struct {
	Direction  entity.CashDirection
	Method     entity.PaymentMethod
	RefType    string
	RefID      string
	SupplierID string
	From       *time.Time
	To         *time.Time
	// Timezone (a Mongo offset such as "+03:30") sets which calendar day a
	// transaction falls on in SummaryByDay; empty means UTC.
	Timezone string
}

type ListCashTransactionsRequest struct {
	Filter     CashFilter
	Pagination commonparam.PaginationRequest
}

type ListCashTransactionsResponse struct {
	Transactions []entity.CashTransaction       `json:"transactions"`
	Summary      CashSummary                    `json:"summary"`
	Pagination   commonparam.PaginationResponse `json:"pagination"`
}

// CashSummary is the دخل و خرج of the selected window.
type CashSummary struct {
	TotalIn  uint64 `json:"total_in"`
	TotalOut uint64 `json:"total_out"`
	Balance  int64  `json:"balance"`
}

type CashMethodRow struct {
	Method   entity.PaymentMethod `json:"method"`
	TotalIn  uint64               `json:"total_in"`
	TotalOut uint64               `json:"total_out"`
	Balance  int64                `json:"balance"`
}

type CashDayRow struct {
	Day      string `json:"day"`
	TotalIn  uint64 `json:"total_in"`
	TotalOut uint64 `json:"total_out"`
	Balance  int64  `json:"balance"`
}

type CashReportRequest struct {
	Filter CashFilter
}

type CashReportResponse struct {
	Summary  CashSummary     `json:"summary"`
	ByMethod []CashMethodRow `json:"by_method"`
	ByDay    []CashDayRow    `json:"by_day"`
}

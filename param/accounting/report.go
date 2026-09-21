package accountingparam

import (
	"time"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type ProfitLossRequest struct {
	From *time.Time
	To   *time.Time
}

// ProfitLossResponse is the صورت سود و زیان: sales, what those sales cost,
// what came back, and what the shop spent to keep the lights on.
type ProfitLossResponse struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`

	UnitsSold  uint64 `json:"units_sold"`
	Revenue    uint64 `json:"revenue"`     // فروش ناخالص
	Returns    uint64 `json:"returns"`     // برگشت از فروش (مبلغ مسترد)
	NetRevenue uint64 `json:"net_revenue"` // فروش خالص

	COGS         uint64 `json:"cogs"`          // بهای تمام‌شده کالای فروش‌رفته
	ReturnedCost uint64 `json:"returned_cost"` // بهای تمام‌شده کالای برگشتی
	NetCOGS      uint64 `json:"net_cogs"`

	GrossProfit int64 `json:"gross_profit"` // سود ناخالص

	Expenses           uint64               `json:"expenses"` // هزینه‌های عملیاتی
	ExpenseBreakdown   []ExpenseCategoryRow `json:"expense_breakdown,omitempty"`
	NetProfit          int64                `json:"net_profit"` // سود خالص
	GrossMarginPercent float64              `json:"gross_margin_percent"`
	NetMarginPercent   float64              `json:"net_margin_percent"`
}

// SupplierLedgerRow is one supplier's account: bought, paid, still owed.
type SupplierLedgerRow struct {
	SupplierID     string     `json:"supplier_id"`
	SupplierName   string     `json:"supplier_name,omitempty"`
	InvoiceCount   int64      `json:"invoice_count"`
	TotalPurchased uint64     `json:"total_purchased"`
	TotalPaid      uint64     `json:"total_paid"`
	Balance        int64      `json:"balance"` // بدهی به تأمین‌کننده (منفی = پیش‌پرداخت)
	LastPurchaseAt *time.Time `json:"last_purchase_at,omitempty"`
}

type SupplierLedgerRequest struct {
	SupplierID string
	OnlyDebt   bool
	From       *time.Time
	To         *time.Time
	Pagination commonparam.PaginationRequest
}

type SupplierLedgerResponse struct {
	Rows           []SupplierLedgerRow `json:"rows"`
	TotalDebt      int64               `json:"total_debt"`
	TotalPurchased uint64              `json:"total_purchased"`
	TotalPaid      uint64              `json:"total_paid"`
}

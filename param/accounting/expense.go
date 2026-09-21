package accountingparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type CreateExpenseRequest struct {
	Title         string                 `json:"title"`
	Category      entity.ExpenseCategory `json:"category"`
	Amount        uint64                 `json:"amount"`
	SupplierID    string                 `json:"supplier_id"`
	PaymentMethod entity.PaymentMethod   `json:"payment_method"`
	IsPaid        *bool                  `json:"is_paid"`
	SpentAt       *time.Time             `json:"spent_at"`
	Note          string                 `json:"note"`
	Attachment    string                 `json:"attachment"`
}

type CreateExpenseResponse struct {
	Expense entity.Expense `json:"expense"`
}

type UpdateExpenseRequest struct {
	ID            string                  `json:"-"`
	Title         *string                 `json:"title"`
	Category      *entity.ExpenseCategory `json:"category"`
	Amount        *uint64                 `json:"amount"`
	SupplierID    *string                 `json:"supplier_id"`
	PaymentMethod *entity.PaymentMethod   `json:"payment_method"`
	IsPaid        *bool                   `json:"is_paid"`
	SpentAt       *time.Time              `json:"spent_at"`
	Note          *string                 `json:"note"`
	Attachment    *string                 `json:"attachment"`
}

type UpdateExpenseResponse struct {
	Expense entity.Expense `json:"expense"`
}

// ExpenseFilter narrows the expense list and every expense-based report.
type ExpenseFilter struct {
	Category   entity.ExpenseCategory
	SupplierID string
	IsPaid     *bool
	Search     string
	From       *time.Time
	To         *time.Time
}

type ListExpensesRequest struct {
	Filter     ExpenseFilter
	Pagination commonparam.PaginationRequest
}

type ListExpensesResponse struct {
	Expenses    []entity.Expense               `json:"expenses"`
	TotalAmount uint64                         `json:"total_amount"`
	Pagination  commonparam.PaginationResponse `json:"pagination"`
}

type GetExpenseResponse struct {
	Expense entity.Expense `json:"expense"`
}

// ExpenseCategoryRow is one bucket of the "کجا خرج شد" breakdown.
type ExpenseCategoryRow struct {
	Category entity.ExpenseCategory `json:"category"`
	Amount   uint64                 `json:"amount"`
	Count    int64                  `json:"count"`
}

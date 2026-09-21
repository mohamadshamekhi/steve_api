package accountingservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	salesreturnparam "github.com/abolfazlnorzad/steve/param/salesreturn"
)

// ExpenseRepository persists operating costs.
type ExpenseRepository interface {
	Create(ctx context.Context, e entity.Expense) error
	FindByID(ctx context.Context, id string) (entity.Expense, error)
	Update(ctx context.Context, e entity.Expense) error
	Delete(ctx context.Context, id string) error
	All(ctx context.Context, f accountingparam.ExpenseFilter, p commonparam.PaginationRequest) ([]entity.Expense, int64, error)
	AllForExport(ctx context.Context, f accountingparam.ExpenseFilter) ([]entity.Expense, error)
	Total(ctx context.Context, f accountingparam.ExpenseFilter) (uint64, error)
	TotalByCategory(ctx context.Context, f accountingparam.ExpenseFilter) ([]accountingparam.ExpenseCategoryRow, error)
}

// CashRepository is the append-only cashbox ledger.
type CashRepository interface {
	Create(ctx context.Context, t entity.CashTransaction) error
	All(ctx context.Context, f accountingparam.CashFilter, p commonparam.PaginationRequest) ([]entity.CashTransaction, int64, error)
	AllForExport(ctx context.Context, f accountingparam.CashFilter) ([]entity.CashTransaction, error)
	Summary(ctx context.Context, f accountingparam.CashFilter) (accountingparam.CashSummary, error)
	SummaryByMethod(ctx context.Context, f accountingparam.CashFilter) ([]accountingparam.CashMethodRow, error)
	SummaryByDay(ctx context.Context, f accountingparam.CashFilter) ([]accountingparam.CashDayRow, error)
	FindByID(ctx context.Context, id string) (bool, error)
	Delete(ctx context.Context, id string) error
}

// ProfitProvider is the sales side of the profit and loss statement, read
// from the stock ledger (revenue and cost of goods sold).
type ProfitProvider interface {
	ProfitReport(ctx context.Context, req inventoryparam.ProfitReportRequest) (inventoryparam.ProfitReportResponse, error)
}

// ReturnProvider supplies the «برگشت از فروش» totals.
type ReturnProvider interface {
	Summary(ctx context.Context, f salesreturnparam.SalesReturnFilter) (salesreturnparam.SalesReturnSummary, error)
}

// PurchaseLedgerProvider supplies what the shop owes its suppliers.
type PurchaseLedgerProvider interface {
	SupplierLedger(ctx context.Context, req accountingparam.SupplierLedgerRequest) ([]accountingparam.SupplierLedgerRow, error)
}

type BackofficeService struct {
	expenseRepo ExpenseRepository
	cashRepo    CashRepository
	profit      ProfitProvider
	returns     ReturnProvider
	purchases   PurchaseLedgerProvider
}

func NewBackofficeService(
	expenseRepo ExpenseRepository,
	cashRepo CashRepository,
	profit ProfitProvider,
	returns ReturnProvider,
	purchases PurchaseLedgerProvider,
) BackofficeService {
	return BackofficeService{
		expenseRepo: expenseRepo,
		cashRepo:    cashRepo,
		profit:      profit,
		returns:     returns,
		purchases:   purchases,
	}
}

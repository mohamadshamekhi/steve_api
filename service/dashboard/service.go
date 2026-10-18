package dashboardservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	dashboardparam "github.com/abolfazlnorzad/steve/param/dashboard"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
)

// ProfitLossProvider is the صورت سود و زیان, the source of the sales card so
// the dashboard and the report page never disagree.
type ProfitLossProvider interface {
	ProfitLoss(ctx context.Context, req accountingparam.ProfitLossRequest) (accountingparam.ProfitLossResponse, error)
	SupplierLedger(ctx context.Context, req accountingparam.SupplierLedgerRequest) (accountingparam.SupplierLedgerResponse, error)
}

// SalesLedger reads the sales side of the stock ledger.
type SalesLedger interface {
	DailySales(ctx context.Context, req inventoryparam.ProfitReportRequest, timezone string) ([]inventoryparam.DailySalesRow, error)
	ProfitBreakdown(ctx context.Context, req inventoryparam.ProfitReportRequest) ([]inventoryparam.ProfitRow, error)
}

type CashRepository interface {
	Summary(ctx context.Context, f accountingparam.CashFilter) (accountingparam.CashSummary, error)
	SummaryByMethod(ctx context.Context, f accountingparam.CashFilter) ([]accountingparam.CashMethodRow, error)
	SummaryByDay(ctx context.Context, f accountingparam.CashFilter) ([]accountingparam.CashDayRow, error)
}

type InventoryProvider interface {
	InventoryValue(ctx context.Context, req inventoryparam.InventoryValueRequest) (inventoryparam.InventoryValueResponse, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]entity.Product, error)
}

type OrderRepository interface {
	CountByStatus(ctx context.Context, from, to time.Time) ([]dashboardparam.OrderStatusCount, error)
	Recent(ctx context.Context, limit int64) ([]entity.Order, error)
}

// CategoryRepository resolves category ids to their names, for the AI
// summary's top-categories breakdown.
type CategoryRepository interface {
	GetCategoriesByIDs(ctx context.Context, ids []string) (map[string]entity.Category, error)
}

type UserRepository interface {
	CountCreatedBetween(ctx context.Context, from, to time.Time) (int64, error)
	GetUsersByIDs(ctx context.Context, ids []string) ([]entity.User, error)
}

type RepairRepository interface {
	CountOpen(ctx context.Context) (int64, error)
}

type PermissionProvider interface {
	PermissionSet(ctx context.Context, backofficeUserID string) (map[string]struct{}, error)
}

type Service struct {
	accounting  ProfitLossProvider
	sales       SalesLedger
	cash        CashRepository
	inventory   InventoryProvider
	orders      OrderRepository
	users       UserRepository
	repairs     RepairRepository
	permissions PermissionProvider
	categories  CategoryRepository
}

func New(
	accounting ProfitLossProvider,
	sales SalesLedger,
	cash CashRepository,
	inventory InventoryProvider,
	orders OrderRepository,
	users UserRepository,
	repairs RepairRepository,
	permissions PermissionProvider,
	categories CategoryRepository,
) Service {
	return Service{
		accounting:  accounting,
		sales:       sales,
		cash:        cash,
		inventory:   inventory,
		orders:      orders,
		users:       users,
		repairs:     repairs,
		permissions: permissions,
		categories:  categories,
	}
}

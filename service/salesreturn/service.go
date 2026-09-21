package salesreturnservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	salesreturnparam "github.com/abolfazlnorzad/steve/param/salesreturn"
)

// Repository persists customer returns.
type Repository interface {
	Create(ctx context.Context, r entity.SalesReturn) error
	FindByID(ctx context.Context, id string) (entity.SalesReturn, error)
	FindByOrderID(ctx context.Context, orderID string) ([]entity.SalesReturn, error)
	Update(ctx context.Context, r entity.SalesReturn) error
	All(ctx context.Context, f salesreturnparam.SalesReturnFilter, p commonparam.PaginationRequest) ([]entity.SalesReturn, int64, error)
	AllForExport(ctx context.Context, f salesreturnparam.SalesReturnFilter) ([]entity.SalesReturn, error)
	Summary(ctx context.Context, f salesreturnparam.SalesReturnFilter) (salesreturnparam.SalesReturnSummary, error)
}

// OrderProvider reads the order a return is booked against; the order is the
// only source for what was sold, at what price and at what cost.
type OrderProvider interface {
	GetOrderByID(ctx context.Context, orderID string) (entity.Order, error)
}

// InventoryProvider puts returned goods back on the shelf and writes the
// ledger line for them.
type InventoryProvider interface {
	RecordOrderReturn(ctx context.Context, order entity.Order, item entity.OrderItem) error
}

// ProductProvider applies the stock delta for restocked units. The inventory
// ledger only records the movement; raising the quantity is the product
// domain's job.
type ProductProvider interface {
	ApplyStockDelta(ctx context.Context, productID string, variantID string, delta int64, avgCost *uint64) error
}

// CashRecorder books the refund as money leaving the till.
type CashRecorder interface {
	Record(ctx context.Context, t entity.CashTransaction) error
}

type BackofficeService struct {
	repo            Repository
	orderProvider   OrderProvider
	inventory       InventoryProvider
	productProvider ProductProvider
	cash            CashRecorder
}

func NewBackofficeService(
	repo Repository,
	orderProvider OrderProvider,
	inventory InventoryProvider,
	productProvider ProductProvider,
) BackofficeService {
	return BackofficeService{
		repo:            repo,
		orderProvider:   orderProvider,
		inventory:       inventory,
		productProvider: productProvider,
	}
}

// SetCashRecorder wires the cashbox after construction, mirroring how the
// order domain receives its inventory recorder — it keeps the dependency
// one-directional and optional.
func (s *BackofficeService) SetCashRecorder(cash CashRecorder) {
	s.cash = cash
}

package purchaseservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	mongodbpurchase "github.com/abolfazlnorzad/steve/repository/mongo/purchase"
)

// Repository is the purchase-invoice persistence port.
type Repository interface {
	Create(ctx context.Context, p entity.Purchase) error
	FindByID(ctx context.Context, id string) (entity.Purchase, error)
	All(ctx context.Context, f mongodbpurchase.PurchaseFilter, p commonparam.PaginationRequest) ([]entity.Purchase, int64, error)
	Update(ctx context.Context, p entity.Purchase) error
}

// InventoryProvider applies a confirmed invoice to the stock ledger.
type InventoryProvider interface {
	ApplyPurchase(ctx context.Context, p entity.Purchase) error
	ReversePurchase(ctx context.Context, p entity.Purchase) error
}

// ProductProvider resolves catalog names and validates that every invoice
// line points at a real variant.
type ProductProvider interface {
	FindByID(ctx context.Context, id string) (entity.Product, error)
}

// SupplierProvider resolves the supplier the invoice is issued by.
type SupplierProvider interface {
	FindByID(ctx context.Context, id string) (entity.Supplier, error)
}

// CashRecorder books supplier payments as money leaving the till.
type CashRecorder interface {
	Record(ctx context.Context, t entity.CashTransaction) error
}

type BackofficeService struct {
	repo             Repository
	inventory        InventoryProvider
	productProvider  ProductProvider
	supplierProvider SupplierProvider
	cash             CashRecorder
}

func NewBackofficeService(
	repo Repository,
	inventory InventoryProvider,
	productProvider ProductProvider,
	supplierProvider SupplierProvider,
) BackofficeService {
	return BackofficeService{
		repo:             repo,
		inventory:        inventory,
		productProvider:  productProvider,
		supplierProvider: supplierProvider,
	}
}

// SetCashRecorder wires the cashbox after construction, keeping the purchase
// domain free of a hard dependency on accounting.
func (s *BackofficeService) SetCashRecorder(cash CashRecorder) {
	s.cash = cash
}

package inventoryservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
)

// Repository is the append-only stock ledger.
type Repository interface {
	Create(ctx context.Context, m entity.StockMovement) error
	CreateMany(ctx context.Context, movements []entity.StockMovement) error
	All(ctx context.Context, f inventoryparam.MovementFilter, p commonparam.PaginationRequest) ([]entity.StockMovement, int64, error)
	DeleteByRef(ctx context.Context, refType string, refID string) error
	ProfitSummary(ctx context.Context, req inventoryparam.ProfitReportRequest) (inventoryparam.ProfitRow, error)
	ProfitBreakdown(ctx context.Context, req inventoryparam.ProfitReportRequest) ([]inventoryparam.ProfitRow, error)
}

// ProductProvider is the slice of the product repo the inventory domain needs:
// reading a variant's current stock/cost and writing the cached totals back.
type ProductProvider interface {
	FindByID(ctx context.Context, id string) (entity.Product, error)
	ApplyStockDelta(ctx context.Context, productID string, variantID string, delta int64, avgCost *uint64) error
	SetVariantStockConfig(ctx context.Context, req inventoryparam.SetStockConfigRequest) error
	InventoryValue(ctx context.Context, req inventoryparam.InventoryValueRequest) (inventoryparam.InventoryValueResponse, error)
}

type BackofficeService struct {
	repo            Repository
	productProvider ProductProvider
}

func NewBackofficeService(repo Repository, productProvider ProductProvider) BackofficeService {
	return BackofficeService{repo: repo, productProvider: productProvider}
}

// findVariant locates one variant inside a product, returning ok=false when
// the product or the variant is gone (a deleted catalog entry must not break
// an accounting write).
func findVariant(p entity.Product, variantID string) (entity.Variant, bool) {
	for _, v := range p.Variants {
		if v.ID == variantID {
			return v, true
		}
	}
	return entity.Variant{}, false
}

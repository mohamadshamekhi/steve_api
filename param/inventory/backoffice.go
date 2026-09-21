package inventoryparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

// MovementFilter narrows a stock-ledger query. Every field is optional.
type MovementFilter struct {
	ProductID  string
	VariantID  string
	SupplierID string
	Type       entity.StockMovementType
	RefType    string
	RefID      string
	From       *time.Time
	To         *time.Time
}

type ListMovementsRequest struct {
	Filter     MovementFilter
	Pagination commonparam.PaginationRequest
}

type ListMovementsResponse struct {
	Movements  []entity.StockMovement         `json:"movements"`
	Pagination commonparam.PaginationResponse `json:"pagination"`
}

// AdjustStockRequest is a manual inventory correction (انبارگردانی).
// Quantity is signed: positive adds units, negative removes them.
type AdjustStockRequest struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
	Quantity  int64  `json:"quantity"`
	UnitCost  uint64 `json:"unit_cost"`
	Note      string `json:"note"`
}

type AdjustStockResponse struct {
	Movement entity.StockMovement `json:"movement"`
}

// SetStockConfigRequest configures how a variant is stocked and sourced.
// Nil fields are left unchanged.
type SetStockConfigRequest struct {
	ProductID      string            `json:"product_id"`
	VariantID      string            `json:"variant_id"`
	StockMode      *entity.StockMode `json:"stock_mode"`
	AllowBackorder *bool             `json:"allow_backorder"`
	BackorderLimit *uint64           `json:"backorder_limit"`
	SupplierID     *string           `json:"supplier_id"`
	SupplyCost     *uint64           `json:"supply_cost"`
	AvgCost        *uint64           `json:"avg_cost"`
}

type SetStockConfigResponse struct{}

// --- reports ---

type ProfitReportRequest struct {
	From    *time.Time
	To      *time.Time
	GroupBy string // day | product | supplier | "" (summary only)
}

// ProfitRow is one bucket of the profit report.
type ProfitRow struct {
	Key         string `json:"key"`             // day / product id / supplier id
	Label       string `json:"label,omitempty"` // human-readable name
	UnitsSold   int64  `json:"units_sold"`
	Revenue     uint64 `json:"revenue"`
	Cost        uint64 `json:"cost"`
	GrossProfit int64  `json:"gross_profit"`
}

type ProfitReportResponse struct {
	Summary ProfitRow   `json:"summary"`
	Rows    []ProfitRow `json:"rows,omitempty"`
}

// InventoryValueRow is one variant's current stock valuation.
type InventoryValueRow struct {
	ProductID   string                    `json:"product_id"`
	ProductName string                    `json:"product_name"`
	VariantID   string                    `json:"variant_id"`
	Attributes  []entity.VariantAttribute `json:"attributes,omitempty"`
	StockMode   entity.StockMode          `json:"stock_mode"`
	Quantity    uint64                    `json:"quantity"`
	AvgCost     uint64                    `json:"avg_cost"`
	Price       uint64                    `json:"price"`
	StockValue  uint64                    `json:"stock_value"`
	// PotentialProfit is what the remaining stock would earn at list price.
	PotentialProfit int64 `json:"potential_profit"`
}

type InventoryValueRequest struct {
	Pagination commonparam.PaginationRequest
	// LowStockOnly limits the listing to tracked variants at or below
	// LowStockThreshold.
	LowStockOnly      bool
	LowStockThreshold uint64
}

type InventoryValueResponse struct {
	Rows            []InventoryValueRow            `json:"rows"`
	TotalStockValue uint64                         `json:"total_stock_value"`
	TotalUnits      uint64                         `json:"total_units"`
	Pagination      commonparam.PaginationResponse `json:"pagination"`
}

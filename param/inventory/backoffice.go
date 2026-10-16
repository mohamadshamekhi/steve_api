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

// StocktakeRequest is a whole-warehouse (or partial) count: the counted
// quantity per variant, from which the server derives each adjustment.
type StocktakeRequest struct {
	Note  string          `json:"note"`
	Items []StocktakeItem `json:"items"`
}

type StocktakeItem struct {
	ProductID       string  `json:"product_id"`
	VariantID       string  `json:"variant_id"`
	CountedQuantity *uint64 `json:"counted_quantity"`
}

type StocktakeResponse struct {
	// StocktakeID is the RefID shared by all movements of this count
	// (GET /backoffice/inventory/movements?ref_type=stocktake&ref_id=...).
	StocktakeID    string `json:"stocktake_id"`
	AdjustedCount  int    `json:"adjusted_count"`
	UnchangedCount int    `json:"unchanged_count"`
	SurplusUnits   uint64 `json:"surplus_units"`
	ShortageUnits  uint64 `json:"shortage_units"`
	// ValueDiff is the signed stock-value change at average cost, in Toman.
	ValueDiff int64                  `json:"value_diff"`
	Movements []entity.StockMovement `json:"movements"`
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
	ProductID   string                    `bson:"product_id" json:"product_id"`
	ProductName string                    `bson:"product_name" json:"product_name"`
	VariantID   string                    `bson:"variant_id" json:"variant_id"`
	Attributes  []entity.VariantAttribute `bson:"attributes,omitempty" json:"attributes,omitempty"`
	StockMode   entity.StockMode          `bson:"stock_mode" json:"stock_mode"`
	Quantity    uint64                    `bson:"quantity" json:"quantity"`
	AvgCost     uint64                    `bson:"avg_cost" json:"avg_cost"`
	Price       uint64                    `bson:"price" json:"price"`
	StockValue  uint64                    `bson:"-" json:"stock_value"`
	// PotentialProfit is what the remaining stock would earn at list price.
	PotentialProfit int64 `bson:"-" json:"potential_profit"`
}

type InventoryValueRequest struct {
	Pagination commonparam.PaginationRequest
	// LowStockOnly limits the listing to variants that still have stock but
	// no more than LowStockThreshold (0 < quantity <= threshold). Sold-out
	// variants are listed by OutOfStockOnly instead.
	LowStockOnly      bool
	LowStockThreshold uint64
	// OutOfStockOnly limits the listing to variants with quantity <= 0.
	OutOfStockOnly bool
	// IncludeSecondHand keeps second-hand products in the low/out-of-stock
	// filters. They are single-unit by nature (quantity 1, then 0), so by
	// default they are left out of those alerts.
	IncludeSecondHand bool
	// Search matches the product name (case-insensitive, substring) or an
	// exact product/variant id.
	Search string
	// StockMode narrows to one stock mode; empty lists every mode except
	// on_demand (which never holds stock).
	StockMode entity.StockMode
}

type InventoryValueResponse struct {
	Rows            []InventoryValueRow            `json:"rows"`
	TotalStockValue uint64                         `json:"total_stock_value"`
	TotalUnits      uint64                         `json:"total_units"`
	Pagination      commonparam.PaginationResponse `json:"pagination"`
}

// DailySalesRow is one calendar day of the sales ledger.
type DailySalesRow struct {
	Day         string `bson:"_id" json:"day"`
	Revenue     uint64 `bson:"revenue" json:"revenue"`
	GrossProfit int64  `bson:"gross_profit" json:"gross_profit"`
	OrderCount  int64  `bson:"order_count" json:"order_count"`
}

// BackfillSalesResponse reports what a sales-ledger backfill wrote.
type BackfillSalesResponse struct {
	Scanned  int `json:"scanned"`
	Recorded int `json:"recorded"`
	Skipped  int `json:"skipped"`
}

package dashboardparam

import commonparam "github.com/abolfazlnorzad/steve/param/common"

// AISummaryRequest is the reporting window plus pagination for the
// (potentially long) per-product inventory list. FromRaw/ToRaw are
// YYYY-MM-DD strings (empty = default 30-day window).
type AISummaryRequest struct {
	FromRaw    string
	ToRaw      string
	Pagination commonparam.PaginationRequest
}

// AISummaryDailySale is one calendar day of aggregated sales. Revenue only —
// no cost/profit.
type AISummaryDailySale struct {
	Day        string `json:"day"`
	OrderCount int64  `json:"order_count"`
	Revenue    uint64 `json:"revenue"`
}

// AISummaryTopProduct is one best-selling product in the window. Revenue
// only — no cost/profit.
type AISummaryTopProduct struct {
	ProductID string `json:"product_id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	UnitsSold int64  `json:"units_sold"`
	Revenue   uint64 `json:"revenue"`
}

// AISummaryTopCategory is one best-selling category in the window, derived
// by summing its products' AISummaryTopProduct rows.
type AISummaryTopCategory struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	UnitsSold  int64  `json:"units_sold"`
	Revenue    uint64 `json:"revenue"`
}

// AISummaryInventoryItem is one product's current stock level (summed
// across its variants). Quantity only — no cost.
type AISummaryInventoryItem struct {
	ProductID string `json:"product_id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Quantity  uint64 `json:"quantity"`
}

// AISummaryResponse is the read-only, PII-free and cost-free sales/inventory
// summary served to the ai_agent role. It must never carry customer
// name/phone, supply_cost, avg_cost, gross_profit, net_profit, net_cogs, or
// any cash/accounting field.
type AISummaryResponse struct {
	Period         Period                         `json:"period"`
	DailySales     []AISummaryDailySale           `json:"daily_sales"`
	OrdersByStatus []OrderStatusCount             `json:"orders_by_status"`
	TopProducts    []AISummaryTopProduct          `json:"top_products"`
	TopCategories  []AISummaryTopCategory         `json:"top_categories"`
	Inventory      []AISummaryInventoryItem       `json:"inventory"`
	InventoryPage  commonparam.PaginationResponse `json:"inventory_page"`
}

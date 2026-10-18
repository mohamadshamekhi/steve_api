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

// AISummaryCancelReasonRow is how many cancelled orders in the window had a
// given CancelReason. "unknown" covers orders cancelled before this field
// existed, which carry no reason.
type AISummaryCancelReasonRow struct {
	Reason string `json:"reason"`
	Count  int64  `json:"count"`
}

// AISummaryCancelPaymentStatusRow is how many cancelled orders in the window
// had a given latest payment status. "no_payment_record" covers orders that
// never got a payment attempt at all.
type AISummaryCancelPaymentStatusRow struct {
	PaymentStatus string `json:"payment_status"`
	Count         int64  `json:"count"`
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

// AISummaryDefinitions documents, inline in the response, exactly what each
// section counts — read this before comparing numbers across sections, they
// are NOT drawn from the same source and will not sum cleanly.
type AISummaryDefinitions struct {
	Currency              string `json:"currency"`
	DailySales            string `json:"daily_sales"`
	TopProductsCategories string `json:"top_products_and_categories"`
	OrdersByStatus        string `json:"orders_by_status"`
	NonCatalogRevenue     string `json:"non_catalog_revenue"`
}

// AISummaryResponse is the read-only, PII-free and cost-free sales/inventory
// summary served to the ai_agent role. It must never carry customer
// name/phone, supply_cost, avg_cost, gross_profit, net_profit, net_cogs, or
// any cash/accounting field.
type AISummaryResponse struct {
	Period      Period               `json:"period"`
	Definitions AISummaryDefinitions `json:"definitions"`
	DailySales  []AISummaryDailySale `json:"daily_sales"`
	// NonCatalogRevenue/Units is the part of DailySales' revenue that came
	// from hand-typed invoice lines (no catalog product_id) — a repair fee,
	// an Apple ID, a part sourced outside the catalog — so it never appears
	// under TopProducts/TopCategories. Compare against daily_sales' total to
	// see how much of a day's revenue is "catalog-attributable".
	NonCatalogRevenue uint64 `json:"non_catalog_revenue"`
	NonCatalogUnits   int64  `json:"non_catalog_units"`

	OrdersByStatus []OrderStatusCount `json:"orders_by_status"`
	// CancelledByReason/CancelledByPaymentStatus only cover orders whose
	// status is cancelled (a subset of OrdersByStatus). Reason "unknown" and
	// payment_status "no_payment_record" are both legitimate historical/edge
	// buckets, not errors.
	CancelledByReason        []AISummaryCancelReasonRow        `json:"cancelled_by_reason"`
	CancelledByPaymentStatus []AISummaryCancelPaymentStatusRow `json:"cancelled_by_payment_status"`

	TopProducts   []AISummaryTopProduct          `json:"top_products"`
	TopCategories []AISummaryTopCategory         `json:"top_categories"`
	Inventory     []AISummaryInventoryItem       `json:"inventory"`
	InventoryPage commonparam.PaginationResponse `json:"inventory_page"`
}

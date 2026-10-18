package dashboardservice

import (
	"context"
	"sort"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	dashboardparam "github.com/abolfazlnorzad/steve/param/dashboard"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"golang.org/x/sync/errgroup"
)

const (
	aiSummaryTopProductsLimit = 20
	aiSummaryTopCategoryLimit = 10
)

// AISummary builds the read-only, aggregated-only sales/inventory snapshot
// served to the ai_agent backoffice role. It must never leak customer PII
// or cost/profit fields — see dashboardparam.AISummaryResponse.
func (s Service) AISummary(ctx context.Context, req dashboardparam.AISummaryRequest) (dashboardparam.AISummaryResponse, error) {
	const op = "dashboardservice.AISummary"

	cur, prev, err := resolvePeriod(req.FromRaw, req.ToRaw, time.Now())
	if err != nil {
		return dashboardparam.AISummaryResponse{}, err
	}

	var (
		byStatus  []dashboardparam.OrderStatusCount
		daily     []inventoryparam.DailySalesRow
		products  []inventoryparam.ProfitRow
		inv       inventoryparam.InventoryValueResponse
		cancelled []entity.Order
	)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() (err error) { byStatus, err = s.orders.CountByStatus(gCtx, cur.from, cur.to); return })
	g.Go(func() (err error) {
		daily, err = s.sales.DailySales(gCtx, inventoryparam.ProfitReportRequest{From: cur.fromPtr(), To: cur.toPtr()}, mongoTimezone)
		return
	})
	g.Go(func() (err error) {
		products, err = s.sales.ProfitBreakdown(gCtx, inventoryparam.ProfitReportRequest{
			From: cur.fromPtr(), To: cur.toPtr(), GroupBy: "product",
		})
		return
	})
	g.Go(func() (err error) {
		inv, err = s.inventory.InventoryValue(gCtx, inventoryparam.InventoryValueRequest{Pagination: req.Pagination})
		return
	})
	g.Go(func() (err error) {
		cancelled, err = s.cancelled.CancelledOrdersBetween(gCtx, cur.from, cur.to)
		return
	})
	if err := g.Wait(); err != nil {
		return dashboardparam.AISummaryResponse{}, richerror.New(op).WithWarpError(err)
	}

	resp := dashboardparam.AISummaryResponse{
		Period: dashboardparam.Period{
			From: cur.fromDay(), To: cur.toDay(),
			PrevFrom: prev.fromDay(), PrevTo: prev.toDay(),
		},
		Definitions: dashboardparam.AISummaryDefinitions{
			Currency: "toman",
			DailySales: "Sales-ledger entries: written when an order is registered " +
				"(before payment), excluded from this query if the order's CURRENT " +
				"status is pending_payment or cancelled. A cancelled order's earlier " +
				"entry is not deleted, only excluded live — so this reflects " +
				"confirmed/non-cancelled orders only, not all registrations.",
			TopProductsCategories: "Same ledger/exclusion rule as daily_sales, grouped " +
				"by product/category. Ranked by revenue. See non_catalog_revenue for the " +
				"part of daily_sales' total these two lists cannot show.",
			OrdersByStatus: "Every order created in the window, grouped by its CURRENT " +
				"status — includes pending_payment and cancelled, unlike daily_sales/" +
				"top_products/top_categories above. This is why its total order count " +
				"is usually much higher.",
			NonCatalogRevenue: "Revenue from hand-typed invoice lines with no catalog " +
				"product_id (e.g. a repair fee, a non-catalog part) — counted in " +
				"daily_sales but not attributable to any product/category.",
		},
		OrdersByStatus: byStatus,
	}

	resp.DailySales = make([]dashboardparam.AISummaryDailySale, 0, len(daily))
	for _, r := range daily {
		resp.DailySales = append(resp.DailySales, dashboardparam.AISummaryDailySale{
			Day: r.Day, OrderCount: r.OrderCount, Revenue: r.Revenue,
		})
	}
	sort.Slice(resp.DailySales, func(i, j int) bool { return resp.DailySales[i].Day < resp.DailySales[j].Day })

	// products is already sorted by gross profit; re-sort by revenue since
	// that is the only money figure this endpoint is allowed to rank by.
	sort.Slice(products, func(i, j int) bool { return products[i].Revenue > products[j].Revenue })

	saleProductIDs := make([]string, 0, len(products))
	for _, p := range products {
		if p.Key != "" {
			saleProductIDs = append(saleProductIDs, p.Key)
		}
	}
	saleProducts, err := s.inventory.GetProductsByIDs(ctx, saleProductIDs)
	if err != nil {
		return dashboardparam.AISummaryResponse{}, richerror.New(op).WithWarpError(err)
	}
	slugByID := make(map[string]string, len(saleProducts))
	categoryByID := make(map[string]string, len(saleProducts))
	for _, p := range saleProducts {
		slugByID[p.ID] = p.Slug
		categoryByID[p.ID] = p.CategoryID
	}

	catAgg := make(map[string]*dashboardparam.AISummaryTopCategory)
	for _, p := range products {
		if p.Key == "" {
			// Hand-typed invoice lines have no product id: their revenue is
			// in daily_sales but can't be attributed to a product/category.
			resp.NonCatalogRevenue += p.Revenue
			resp.NonCatalogUnits += p.UnitsSold
			continue
		}
		if len(resp.TopProducts) < aiSummaryTopProductsLimit {
			resp.TopProducts = append(resp.TopProducts, dashboardparam.AISummaryTopProduct{
				ProductID: p.Key, Slug: slugByID[p.Key], Name: p.Label,
				UnitsSold: p.UnitsSold, Revenue: p.Revenue,
			})
		}
		if catID := categoryByID[p.Key]; catID != "" {
			c, ok := catAgg[catID]
			if !ok {
				c = &dashboardparam.AISummaryTopCategory{CategoryID: catID}
				catAgg[catID] = c
			}
			c.UnitsSold += p.UnitsSold
			c.Revenue += p.Revenue
		}
	}

	catIDs := make([]string, 0, len(catAgg))
	for id := range catAgg {
		catIDs = append(catIDs, id)
	}
	cats, err := s.categories.GetCategoriesByIDs(ctx, catIDs)
	if err != nil {
		return dashboardparam.AISummaryResponse{}, richerror.New(op).WithWarpError(err)
	}
	resp.TopCategories = make([]dashboardparam.AISummaryTopCategory, 0, len(catAgg))
	for id, c := range catAgg {
		if cat, ok := cats[id]; ok {
			c.Name = cat.Name
		}
		resp.TopCategories = append(resp.TopCategories, *c)
	}
	sort.Slice(resp.TopCategories, func(i, j int) bool { return resp.TopCategories[i].Revenue > resp.TopCategories[j].Revenue })
	if len(resp.TopCategories) > aiSummaryTopCategoryLimit {
		resp.TopCategories = resp.TopCategories[:aiSummaryTopCategoryLimit]
	}

	// Inventory: InventoryValue returns one row per variant; sum to
	// per-product quantity and resolve slugs for the page's products only.
	invProductIDs := make([]string, 0, len(inv.Rows))
	qtyByProduct := make(map[string]uint64, len(inv.Rows))
	nameByProduct := make(map[string]string, len(inv.Rows))
	for _, r := range inv.Rows {
		if _, ok := qtyByProduct[r.ProductID]; !ok {
			invProductIDs = append(invProductIDs, r.ProductID)
		}
		qtyByProduct[r.ProductID] += r.Quantity
		nameByProduct[r.ProductID] = r.ProductName
	}
	invProducts, err := s.inventory.GetProductsByIDs(ctx, invProductIDs)
	if err != nil {
		return dashboardparam.AISummaryResponse{}, richerror.New(op).WithWarpError(err)
	}
	invSlugByID := make(map[string]string, len(invProducts))
	for _, p := range invProducts {
		invSlugByID[p.ID] = p.Slug
	}

	resp.Inventory = make([]dashboardparam.AISummaryInventoryItem, 0, len(invProductIDs))
	for _, pid := range invProductIDs {
		resp.Inventory = append(resp.Inventory, dashboardparam.AISummaryInventoryItem{
			ProductID: pid, Slug: invSlugByID[pid], Name: nameByProduct[pid], Quantity: qtyByProduct[pid],
		})
	}
	resp.InventoryPage = inv.Pagination

	cancelOrderIDs := make([]string, 0, len(cancelled))
	reasonCounts := make(map[entity.CancelReason]int64, len(cancelled))
	for _, o := range cancelled {
		cancelOrderIDs = append(cancelOrderIDs, o.ID)
		reason := o.CancelReason
		if reason == "" {
			reason = "unknown"
		}
		reasonCounts[reason]++
	}
	resp.CancelledByReason = make([]dashboardparam.AISummaryCancelReasonRow, 0, len(reasonCounts))
	for reason, count := range reasonCounts {
		resp.CancelledByReason = append(resp.CancelledByReason, dashboardparam.AISummaryCancelReasonRow{
			Reason: string(reason), Count: count,
		})
	}
	sort.Slice(resp.CancelledByReason, func(i, j int) bool { return resp.CancelledByReason[i].Count > resp.CancelledByReason[j].Count })

	paymentStatusByOrder, err := s.payments.GetLatestStatusByOrderIDs(ctx, cancelOrderIDs)
	if err != nil {
		return dashboardparam.AISummaryResponse{}, richerror.New(op).WithWarpError(err)
	}
	paymentCounts := make(map[string]int64, len(cancelOrderIDs))
	for _, id := range cancelOrderIDs {
		status, ok := paymentStatusByOrder[id]
		key := "no_payment_record"
		if ok {
			key = string(status)
		}
		paymentCounts[key]++
	}
	resp.CancelledByPaymentStatus = make([]dashboardparam.AISummaryCancelPaymentStatusRow, 0, len(paymentCounts))
	for status, count := range paymentCounts {
		resp.CancelledByPaymentStatus = append(resp.CancelledByPaymentStatus, dashboardparam.AISummaryCancelPaymentStatusRow{
			PaymentStatus: status, Count: count,
		})
	}
	sort.Slice(resp.CancelledByPaymentStatus, func(i, j int) bool {
		return resp.CancelledByPaymentStatus[i].Count > resp.CancelledByPaymentStatus[j].Count
	})

	return resp, nil
}

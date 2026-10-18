package dashboardservice

import (
	"context"
	"sort"
	"time"

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

	cur, _, err := resolvePeriod(req.FromRaw, req.ToRaw, time.Now())
	if err != nil {
		return dashboardparam.AISummaryResponse{}, err
	}

	var (
		byStatus []dashboardparam.OrderStatusCount
		daily    []inventoryparam.DailySalesRow
		products []inventoryparam.ProfitRow
		inv      inventoryparam.InventoryValueResponse
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
	if err := g.Wait(); err != nil {
		return dashboardparam.AISummaryResponse{}, richerror.New(op).WithWarpError(err)
	}

	resp := dashboardparam.AISummaryResponse{
		Period:         dashboardparam.Period{From: cur.fromDay(), To: cur.toDay()},
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
			continue // hand-typed invoice lines have no product id
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

	return resp, nil
}

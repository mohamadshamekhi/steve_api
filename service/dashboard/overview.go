package dashboardservice

import (
	"context"
	"sort"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	dashboardparam "github.com/abolfazlnorzad/steve/param/dashboard"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"golang.org/x/sync/errgroup"
)

const (
	// listLimit caps every "top N" list on the dashboard.
	listLimit = 5
	// lowStockThreshold matches the default of GET /backoffice/inventory/stock.
	lowStockThreshold = 3
)

// Overview builds the backoffice home page. Each section is computed only
// when the caller holds its permission and is left null otherwise; the
// sections are independent, so they run concurrently.
func (s Service) Overview(ctx context.Context, req dashboardparam.OverviewRequest) (dashboardparam.OverviewResponse, error) {
	const op = "dashboardservice.Overview"

	cur, prev, err := resolvePeriod(req.FromRaw, req.ToRaw, time.Now())
	if err != nil {
		return dashboardparam.OverviewResponse{}, err
	}

	perms, err := s.permissions.PermissionSet(ctx, req.BackofficeUserID)
	if err != nil {
		return dashboardparam.OverviewResponse{}, richerror.New(op).WithWarpError(err)
	}
	has := func(k string) bool { _, ok := perms[k]; return ok }

	resp := dashboardparam.OverviewResponse{
		Period: dashboardparam.Period{
			From: cur.fromDay(), To: cur.toDay(),
			PrevFrom: prev.fromDay(), PrevTo: prev.toDay(),
		},
	}

	// Every section writes only its own fields of resp, so the goroutines
	// never touch the same memory.
	g, gCtx := errgroup.WithContext(ctx)
	if has(key.ViewProfitLossReport) {
		g.Go(func() error { return s.fillSales(gCtx, &resp, cur, prev) })
	}
	if has(key.ViewCashReport) {
		g.Go(func() error { return s.fillCash(gCtx, &resp, cur) })
	}
	if has(key.ViewInventory) {
		g.Go(func() error { return s.fillInventory(gCtx, &resp) })
	}
	if has(key.ViewSupplierLedger) {
		g.Go(func() error { return s.fillSuppliers(gCtx, &resp) })
	}
	if has(key.ListOrders) {
		g.Go(func() error { return s.fillOrders(gCtx, &resp, cur) })
	}
	if has(key.ListUsers) {
		g.Go(func() error { return s.fillCustomers(gCtx, &resp, cur, prev) })
	}
	if has(key.ListRepairRequests) {
		g.Go(func() error { return s.fillRepairs(gCtx, &resp) })
	}
	if err := g.Wait(); err != nil {
		return dashboardparam.OverviewResponse{}, richerror.New(op).WithWarpError(err)
	}
	return resp, nil
}

// fillSales fills sales, trend, top_products and expense_breakdown, all
// read from the same sources as GET /backoffice/report/profit-loss.
func (s Service) fillSales(ctx context.Context, resp *dashboardparam.OverviewResponse, cur, prev window) error {
	var (
		pl, plPrev       accountingparam.ProfitLossResponse
		daily, dailyPrev []inventoryparam.DailySalesRow
		products         []inventoryparam.ProfitRow
	)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		pl, err = s.accounting.ProfitLoss(gCtx, accountingparam.ProfitLossRequest{From: cur.fromPtr(), To: cur.toPtr()})
		return err
	})
	g.Go(func() (err error) {
		plPrev, err = s.accounting.ProfitLoss(gCtx, accountingparam.ProfitLossRequest{From: prev.fromPtr(), To: prev.toPtr()})
		return err
	})
	g.Go(func() (err error) {
		daily, err = s.sales.DailySales(gCtx, inventoryparam.ProfitReportRequest{From: cur.fromPtr(), To: cur.toPtr()}, mongoTimezone)
		return err
	})
	g.Go(func() (err error) {
		dailyPrev, err = s.sales.DailySales(gCtx, inventoryparam.ProfitReportRequest{From: prev.fromPtr(), To: prev.toPtr()}, mongoTimezone)
		return err
	})
	g.Go(func() (err error) {
		products, err = s.sales.ProfitBreakdown(gCtx, inventoryparam.ProfitReportRequest{
			From: cur.fromPtr(), To: cur.toPtr(), GroupBy: "product",
		})
		return err
	})
	if err := g.Wait(); err != nil {
		return err
	}

	sales := &dashboardparam.Sales{
		NetRevenue:         pl.NetRevenue,
		NetRevenuePrev:     plPrev.NetRevenue,
		NetProfit:          pl.NetProfit,
		NetProfitPrev:      plPrev.NetProfit,
		GrossProfit:        pl.GrossProfit,
		Revenue:            pl.Revenue,
		Returns:            pl.Returns,
		NetCOGS:            pl.NetCOGS,
		Expenses:           pl.Expenses,
		GrossMarginPercent: pl.GrossMarginPercent,
		NetMarginPercent:   pl.NetMarginPercent,
		UnitsSold:          pl.UnitsSold,
		OrderCount:         sumOrders(daily),
		OrderCountPrev:     sumOrders(dailyPrev),
	}
	if sales.OrderCount > 0 {
		sales.AvgOrderValue = pl.Revenue / uint64(sales.OrderCount)
	}
	resp.Sales = sales

	byDay := make(map[string]inventoryparam.DailySalesRow, len(daily))
	for _, r := range daily {
		byDay[r.Day] = r
	}
	resp.Trend = make([]dashboardparam.TrendRow, 0, cur.days)
	for _, day := range cur.dayKeys() {
		r := byDay[day]
		resp.Trend = append(resp.Trend, dashboardparam.TrendRow{
			Day: day, Revenue: r.Revenue, GrossProfit: r.GrossProfit, OrderCount: r.OrderCount,
		})
	}

	// ProfitBreakdown is already ordered by gross profit. Hand-typed counter
	// lines have no product id and are not a product to rank.
	resp.TopProducts = make([]dashboardparam.TopProduct, 0, listLimit)
	for _, p := range products {
		if p.Key == "" {
			continue
		}
		resp.TopProducts = append(resp.TopProducts, dashboardparam.TopProduct{
			ProductID: p.Key, ProductName: p.Label,
			UnitsSold: p.UnitsSold, Revenue: p.Revenue, GrossProfit: p.GrossProfit,
		})
		if len(resp.TopProducts) == listLimit {
			break
		}
	}

	resp.ExpenseBreakdown = pl.ExpenseBreakdown
	if resp.ExpenseBreakdown == nil {
		resp.ExpenseBreakdown = []accountingparam.ExpenseCategoryRow{}
	}
	return nil
}

func sumOrders(rows []inventoryparam.DailySalesRow) int64 {
	var n int64
	for _, r := range rows {
		n += r.OrderCount
	}
	return n
}

func (s Service) fillCash(ctx context.Context, resp *dashboardparam.OverviewResponse, cur window) error {
	f := accountingparam.CashFilter{From: cur.fromPtr(), To: cur.toPtr(), Timezone: mongoTimezone}

	var (
		summary  accountingparam.CashSummary
		byDay    []accountingparam.CashDayRow
		byMethod []accountingparam.CashMethodRow
	)
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() (err error) { summary, err = s.cash.Summary(gCtx, f); return err })
	g.Go(func() (err error) { byDay, err = s.cash.SummaryByDay(gCtx, f); return err })
	g.Go(func() (err error) { byMethod, err = s.cash.SummaryByMethod(gCtx, f); return err })
	if err := g.Wait(); err != nil {
		return err
	}

	cash := &dashboardparam.Cash{
		TotalIn:  summary.TotalIn,
		TotalOut: summary.TotalOut,
		Balance:  summary.Balance,
		ByDay:    make([]dashboardparam.CashDay, 0, cur.days),
		ByMethod: make([]dashboardparam.CashMethod, 0, len(byMethod)),
	}

	dayRows := make(map[string]accountingparam.CashDayRow, len(byDay))
	for _, r := range byDay {
		dayRows[r.Day] = r
	}
	for _, day := range cur.dayKeys() {
		r := dayRows[day]
		cash.ByDay = append(cash.ByDay, dashboardparam.CashDay{Day: day, TotalIn: r.TotalIn, TotalOut: r.TotalOut})
	}

	for _, m := range byMethod {
		if m.TotalIn == 0 {
			continue
		}
		cash.ByMethod = append(cash.ByMethod, dashboardparam.CashMethod{Method: m.Method, TotalIn: m.TotalIn})
	}
	sort.Slice(cash.ByMethod, func(i, j int) bool { return cash.ByMethod[i].TotalIn > cash.ByMethod[j].TotalIn })

	resp.Cash = cash
	return nil
}

func (s Service) fillInventory(ctx context.Context, resp *dashboardparam.OverviewResponse) error {
	one := commonparam.PaginationRequest{PageNumber: 1, PageSize: 1}

	var all, low, out inventoryparam.InventoryValueResponse
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		all, err = s.inventory.InventoryValue(gCtx, inventoryparam.InventoryValueRequest{Pagination: one})
		return err
	})
	g.Go(func() (err error) {
		low, err = s.inventory.InventoryValue(gCtx, inventoryparam.InventoryValueRequest{
			// 5 is not an allowed page size; fetch 10 and keep the first 5.
			Pagination:        commonparam.PaginationRequest{PageNumber: 1, PageSize: 10},
			LowStockOnly:      true,
			LowStockThreshold: lowStockThreshold,
		})
		return err
	})
	g.Go(func() (err error) {
		out, err = s.inventory.InventoryValue(gCtx, inventoryparam.InventoryValueRequest{Pagination: one, OutOfStockOnly: true})
		return err
	})
	if err := g.Wait(); err != nil {
		return err
	}

	inv := &dashboardparam.Inventory{
		TotalStockValue:   all.TotalStockValue,
		TotalUnits:        all.TotalUnits,
		LowStockCount:     low.Pagination.Total,
		OutOfStockCount:   out.Pagination.Total,
		LowStockThreshold: lowStockThreshold,
		LowStock:          make([]dashboardparam.LowStockItem, 0, len(low.Rows)),
	}
	for _, r := range low.Rows {
		if len(inv.LowStock) == listLimit {
			break
		}
		attrs := r.Attributes
		if attrs == nil {
			attrs = []entity.VariantAttribute{}
		}
		inv.LowStock = append(inv.LowStock, dashboardparam.LowStockItem{
			ProductID: r.ProductID, VariantID: r.VariantID, ProductName: r.ProductName,
			Attributes: attrs, Quantity: r.Quantity,
		})
	}
	resp.Inventory = inv
	return nil
}

func (s Service) fillSuppliers(ctx context.Context, resp *dashboardparam.OverviewResponse) error {
	ledger, err := s.accounting.SupplierLedger(ctx, accountingparam.SupplierLedgerRequest{OnlyDebt: true})
	if err != nil {
		return err
	}

	rows := ledger.Rows
	sort.Slice(rows, func(i, j int) bool { return rows[i].Balance > rows[j].Balance })

	sup := &dashboardparam.Suppliers{
		TotalDebt:  ledger.TotalDebt,
		TopDebtors: make([]dashboardparam.SupplierDebtor, 0, listLimit),
	}
	for _, r := range rows {
		if len(sup.TopDebtors) == listLimit {
			break
		}
		sup.TopDebtors = append(sup.TopDebtors, dashboardparam.SupplierDebtor{
			SupplierID: r.SupplierID, SupplierName: r.SupplierName, Balance: r.Balance,
		})
	}
	resp.Suppliers = sup
	return nil
}

func (s Service) fillOrders(ctx context.Context, resp *dashboardparam.OverviewResponse, cur window) error {
	byStatus, err := s.orders.CountByStatus(ctx, cur.from, cur.to)
	if err != nil {
		return err
	}
	orders, err := s.orders.Recent(ctx, listLimit)
	if err != nil {
		return err
	}

	userIDs := make([]string, 0, len(orders))
	for _, o := range orders {
		if o.UserID != "" {
			userIDs = append(userIDs, o.UserID)
		}
	}
	users, err := s.users.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return err
	}
	names := make(map[string][2]string, len(users))
	for _, u := range users {
		names[u.ID] = [2]string{u.Name, u.PhoneNumber}
	}

	recent := make([]dashboardparam.RecentOrder, 0, len(orders))
	for _, o := range orders {
		u := names[o.UserID]
		recent = append(recent, dashboardparam.RecentOrder{
			ID: o.ID, UserName: u[0], UserPhoneNumber: u[1],
			FinalAmount: o.FinalAmount, Status: o.Status, CreatedAt: o.CreatedAt,
		})
	}

	resp.OrdersByStatus = byStatus
	resp.RecentOrders = recent
	return nil
}

func (s Service) fillCustomers(ctx context.Context, resp *dashboardparam.OverviewResponse, cur, prev window) error {
	n, err := s.users.CountCreatedBetween(ctx, cur.from, cur.to)
	if err != nil {
		return err
	}
	nPrev, err := s.users.CountCreatedBetween(ctx, prev.from, prev.to)
	if err != nil {
		return err
	}
	resp.Customers = &dashboardparam.Customers{NewCount: n, NewCountPrev: nPrev}
	return nil
}

func (s Service) fillRepairs(ctx context.Context, resp *dashboardparam.OverviewResponse) error {
	n, err := s.repairs.CountOpen(ctx)
	if err != nil {
		return err
	}
	resp.Repairs = &dashboardparam.Repairs{OpenCount: n}
	return nil
}

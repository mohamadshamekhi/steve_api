package accountingservice

import (
	"context"

	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	salesreturnparam "github.com/abolfazlnorzad/steve/param/salesreturn"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// ProfitLoss builds the صورت سود و زیان for a window:
//
//	فروش ناخالص − برگشت از فروش            = فروش خالص
//	بهای تمام‌شده − بهای کالای برگشتی        = بهای تمام‌شده خالص
//	فروش خالص − بهای تمام‌شده خالص           = سود ناخالص
//	سود ناخالص − هزینه‌های عملیاتی          = سود خالص
//
// Revenue and cost come from the stock ledger (accrual), not from the
// cashbox: goods sold on credit are still a sale, and money received for a
// previous month is not this month's revenue.
func (s BackofficeService) ProfitLoss(ctx context.Context, req accountingparam.ProfitLossRequest) (accountingparam.ProfitLossResponse, error) {
	const op = "accountingservice.ProfitLoss"

	sales, err := s.profit.ProfitReport(ctx, inventoryparam.ProfitReportRequest{From: req.From, To: req.To})
	if err != nil {
		return accountingparam.ProfitLossResponse{}, richerror.New(op).WithWarpError(err)
	}

	returns, err := s.returns.Summary(ctx, salesreturnparam.SalesReturnFilter{From: req.From, To: req.To})
	if err != nil {
		return accountingparam.ProfitLossResponse{}, richerror.New(op).WithWarpError(err)
	}

	expenseFilter := accountingparam.ExpenseFilter{From: req.From, To: req.To}
	expenses, err := s.expenseRepo.Total(ctx, expenseFilter)
	if err != nil {
		return accountingparam.ProfitLossResponse{}, richerror.New(op).WithWarpError(err)
	}
	breakdown, err := s.expenseRepo.TotalByCategory(ctx, expenseFilter)
	if err != nil {
		return accountingparam.ProfitLossResponse{}, richerror.New(op).WithWarpError(err)
	}

	resp := accountingparam.ProfitLossResponse{
		From:             req.From,
		To:               req.To,
		UnitsSold:        uint64(max64(sales.Summary.UnitsSold, 0)),
		Revenue:          sales.Summary.Revenue,
		Returns:          returns.RefundAmount,
		COGS:             sales.Summary.Cost,
		ReturnedCost:     returns.ReturnedCost,
		Expenses:         expenses,
		ExpenseBreakdown: breakdown,
	}

	// Subtractions are floored at zero: a window where more was refunded than
	// sold (a return booked against last month's order) must not wrap around
	// an unsigned total.
	resp.NetRevenue = subFloor(resp.Revenue, resp.Returns)
	resp.NetCOGS = subFloor(resp.COGS, resp.ReturnedCost)
	resp.GrossProfit = int64(resp.NetRevenue) - int64(resp.NetCOGS)
	resp.NetProfit = resp.GrossProfit - int64(resp.Expenses)

	if resp.NetRevenue > 0 {
		resp.GrossMarginPercent = percent(resp.GrossProfit, resp.NetRevenue)
		resp.NetMarginPercent = percent(resp.NetProfit, resp.NetRevenue)
	}
	return resp, nil
}

// SupplierLedger reports what the shop owes each supplier: purchased, paid
// and the remaining balance.
func (s BackofficeService) SupplierLedger(ctx context.Context, req accountingparam.SupplierLedgerRequest) (accountingparam.SupplierLedgerResponse, error) {
	const op = "accountingservice.SupplierLedger"

	rows, err := s.purchases.SupplierLedger(ctx, req)
	if err != nil {
		return accountingparam.SupplierLedgerResponse{}, richerror.New(op).WithWarpError(err)
	}

	resp := accountingparam.SupplierLedgerResponse{Rows: rows}
	for _, r := range rows {
		resp.TotalPurchased += r.TotalPurchased
		resp.TotalPaid += r.TotalPaid
		resp.TotalDebt += r.Balance
	}
	return resp, nil
}

func subFloor(a, b uint64) uint64 {
	if b >= a {
		return 0
	}
	return a - b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// percent returns part/whole as a percentage rounded to two decimals.
func percent(part int64, whole uint64) float64 {
	if whole == 0 {
		return 0
	}
	return float64(int64(part)*10000/int64(whole)) / 100
}

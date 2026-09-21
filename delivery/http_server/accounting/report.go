package accountinghandler

import (
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	"github.com/abolfazlnorzad/steve/pkg/export"
	"github.com/labstack/echo/v4"
)

// ProfitLoss is the صورت سود و زیان over a window.
func (h Handler) ProfitLoss(c echo.Context) error {
	req := accountingparam.ProfitLossRequest{
		From: parseDate(c.QueryParam("from")),
		To:   parseDate(c.QueryParam("to")),
	}

	resp, err := h.backofficeSvc.ProfitLoss(c.Request().Context(), req)
	if err != nil {
		return err
	}

	if wantsCSV(c) {
		rows := [][]string{
			{"فروش ناخالص", export.Uint(resp.Revenue)},
			{"برگشت از فروش", export.Uint(resp.Returns)},
			{"فروش خالص", export.Uint(resp.NetRevenue)},
			{"بهای تمام‌شده کالای فروش‌رفته", export.Uint(resp.COGS)},
			{"بهای تمام‌شده کالای برگشتی", export.Uint(resp.ReturnedCost)},
			{"بهای تمام‌شده خالص", export.Uint(resp.NetCOGS)},
			{"سود ناخالص", export.Int(resp.GrossProfit)},
			{"هزینه‌های عملیاتی", export.Uint(resp.Expenses)},
			{"سود خالص", export.Int(resp.NetProfit)},
			{"حاشیه سود ناخالص (٪)", export.Float(resp.GrossMarginPercent)},
			{"حاشیه سود خالص (٪)", export.Float(resp.NetMarginPercent)},
		}
		// The category breakdown follows the statement so one file answers
		// both «چقدر سود کردم» and «پولم کجا رفت».
		for _, r := range resp.ExpenseBreakdown {
			rows = append(rows, []string{"هزینه: " + string(r.Category), export.Uint(r.Amount)})
		}
		return export.CSV(c, "profit-loss", []string{"عنوان", "مبلغ"}, rows)
	}
	return h.success(c, resp)
}

// SupplierLedger reports what the shop still owes each supplier.
func (h Handler) SupplierLedger(c echo.Context) error {
	req := accountingparam.SupplierLedgerRequest{
		SupplierID: c.QueryParam("supplier_id"),
		OnlyDebt:   c.QueryParam("only_debt") == "true",
		From:       parseDate(c.QueryParam("from")),
		To:         parseDate(c.QueryParam("to")),
		Pagination: pagination(c),
	}

	resp, err := h.backofficeSvc.SupplierLedger(c.Request().Context(), req)
	if err != nil {
		return err
	}

	if wantsCSV(c) {
		rows := make([][]string, 0, len(resp.Rows)+1)
		for _, r := range resp.Rows {
			rows = append(rows, []string{
				r.SupplierName, export.Int(r.InvoiceCount),
				export.Uint(r.TotalPurchased), export.Uint(r.TotalPaid), export.Int(r.Balance),
				dateOrEmpty(r),
			})
		}
		rows = append(rows, []string{"جمع کل", "",
			export.Uint(resp.TotalPurchased), export.Uint(resp.TotalPaid), export.Int(resp.TotalDebt), ""})
		return export.CSV(c, "supplier-ledger",
			[]string{"تأمین‌کننده", "تعداد فاکتور", "مبلغ خرید", "پرداخت‌شده", "بدهی", "آخرین خرید"}, rows)
	}
	return h.success(c, resp)
}

func dateOrEmpty(r accountingparam.SupplierLedgerRow) string {
	if r.LastPurchaseAt == nil {
		return ""
	}
	return export.Date(*r.LastPurchaseAt)
}

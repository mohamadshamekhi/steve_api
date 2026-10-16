package accountinghandler

import (
	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/export"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func cashFilter(c echo.Context) accountingparam.CashFilter {
	return accountingparam.CashFilter{
		Direction:  entity.CashDirection(c.QueryParam("direction")),
		Method:     entity.PaymentMethod(c.QueryParam("method")),
		RefType:    c.QueryParam("ref_type"),
		RefID:      c.QueryParam("ref_id"),
		SupplierID: c.QueryParam("supplier_id"),
		From:       parseDate(c.QueryParam("from")),
		To:         parseDateEnd(c.QueryParam("to")),
	}
}

func (h Handler) ListCashTransactions(c echo.Context) error {
	filter := cashFilter(c)

	if wantsCSV(c) {
		transactions, err := h.backofficeSvc.ExportCashTransactions(c.Request().Context(), filter)
		if err != nil {
			return err
		}
		rows := make([][]string, 0, len(transactions))
		for _, t := range transactions {
			rows = append(rows, []string{
				export.DateTime(t.OccurredAt), string(t.Direction), string(t.Method),
				export.Uint(t.Amount), t.Description, t.RefType, t.RefID,
			})
		}
		return export.CSV(c, "cash-transactions",
			[]string{"تاریخ", "نوع", "روش", "مبلغ", "شرح", "مرجع", "شناسه مرجع"}, rows)
	}

	resp, err := h.backofficeSvc.ListCashTransactions(c.Request().Context(), accountingparam.ListCashTransactionsRequest{
		Filter:     filter,
		Pagination: pagination(c),
	})
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) CreateCashTransaction(c echo.Context) error {
	const op = "accountinghandler.CreateCashTransaction"

	var req accountingparam.CreateCashTransactionRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	if err := h.validator.ValidateCreateCashTransactionRequest(req); err != nil {
		return err
	}

	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.backofficeSvc.CreateCashTransaction(c.Request().Context(), req, cl.EntityID)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) DeleteCashTransaction(c echo.Context) error {
	if err := h.backofficeSvc.DeleteCashTransaction(c.Request().Context(), c.Param("id")); err != nil {
		return err
	}
	return h.success(c, nil)
}

// CashReport is the دخل و خرج view: totals, split by payment method and by day.
func (h Handler) CashReport(c echo.Context) error {
	filter := cashFilter(c)

	resp, err := h.backofficeSvc.CashReport(c.Request().Context(), accountingparam.CashReportRequest{Filter: filter})
	if err != nil {
		return err
	}

	if wantsCSV(c) {
		rows := make([][]string, 0, len(resp.ByDay))
		for _, r := range resp.ByDay {
			rows = append(rows, []string{r.Day, export.Uint(r.TotalIn), export.Uint(r.TotalOut), export.Int(r.Balance)})
		}
		rows = append(rows, []string{"جمع کل",
			export.Uint(resp.Summary.TotalIn), export.Uint(resp.Summary.TotalOut), export.Int(resp.Summary.Balance)})
		return export.CSV(c, "cash-report", []string{"روز", "دریافت", "پرداخت", "مانده"}, rows)
	}
	return h.success(c, resp)
}

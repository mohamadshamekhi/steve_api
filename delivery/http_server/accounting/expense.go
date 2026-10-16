package accountinghandler

import (
	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/export"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func expenseFilter(c echo.Context) accountingparam.ExpenseFilter {
	return accountingparam.ExpenseFilter{
		Category:   entity.ExpenseCategory(c.QueryParam("category")),
		SupplierID: c.QueryParam("supplier_id"),
		IsPaid:     parseBool(c.QueryParam("is_paid")),
		Search:     c.QueryParam("search"),
		From:       parseDate(c.QueryParam("from")),
		To:         parseDateEnd(c.QueryParam("to")),
	}
}

func (h Handler) ListExpenses(c echo.Context) error {
	filter := expenseFilter(c)

	if wantsCSV(c) {
		expenses, err := h.backofficeSvc.ExportExpenses(c.Request().Context(), filter)
		if err != nil {
			return err
		}
		rows := make([][]string, 0, len(expenses))
		for _, e := range expenses {
			rows = append(rows, []string{
				export.Date(e.SpentAt), e.Title, string(e.Category),
				export.Uint(e.Amount), string(e.PaymentMethod), export.Bool(e.IsPaid), e.Note,
			})
		}
		return export.CSV(c, "expenses",
			[]string{"تاریخ", "عنوان", "دسته", "مبلغ", "روش پرداخت", "پرداخت شده", "توضیح"}, rows)
	}

	resp, err := h.backofficeSvc.ListExpenses(c.Request().Context(), accountingparam.ListExpensesRequest{
		Filter:     filter,
		Pagination: pagination(c),
	})
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) GetExpense(c echo.Context) error {
	resp, err := h.backofficeSvc.GetExpense(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) CreateExpense(c echo.Context) error {
	const op = "accountinghandler.CreateExpense"

	var req accountingparam.CreateExpenseRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	if err := h.validator.ValidateCreateExpenseRequest(req); err != nil {
		return err
	}

	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.backofficeSvc.CreateExpense(c.Request().Context(), req, cl.EntityID)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) UpdateExpense(c echo.Context) error {
	const op = "accountinghandler.UpdateExpense"

	var req accountingparam.UpdateExpenseRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	req.ID = c.Param("id")
	if err := h.validator.ValidateUpdateExpenseRequest(req); err != nil {
		return err
	}

	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.backofficeSvc.UpdateExpense(c.Request().Context(), req, cl.EntityID)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) DeleteExpense(c echo.Context) error {
	if err := h.backofficeSvc.DeleteExpense(c.Request().Context(), c.Param("id")); err != nil {
		return err
	}
	return h.success(c, nil)
}

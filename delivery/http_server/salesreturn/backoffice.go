package salesreturnhandler

import (
	salesreturnparam "github.com/abolfazlnorzad/steve/param/salesreturn"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/export"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) List(c echo.Context) error {
	filter := salesreturnparam.SalesReturnFilter{
		OrderID: c.QueryParam("order_id"),
		UserID:  c.QueryParam("user_id"),
		From:    parseDate(c.QueryParam("from")),
		To:      parseDate(c.QueryParam("to")),
	}

	if c.QueryParam("format") == "csv" {
		returns, err := h.backofficeSvc.Export(c.Request().Context(), filter)
		if err != nil {
			return err
		}
		// One row per returned line, so the file can be pivoted by product.
		rows := make([][]string, 0, len(returns))
		for _, r := range returns {
			for _, item := range r.Items {
				rows = append(rows, []string{
					export.DateTime(r.ReturnedAt), r.OrderID, item.ProductName,
					export.Uint(item.Quantity), export.Uint(item.UnitPrice), export.Uint(item.UnitCost),
					export.Bool(item.Restock), export.Uint(r.RefundAmount), r.Reason,
				})
			}
		}
		return export.CSV(c, "sales-returns",
			[]string{"تاریخ", "سفارش", "کالا", "تعداد", "قیمت فروش", "قیمت تمام‌شده", "برگشت به انبار", "مبلغ مسترد", "دلیل"}, rows)
	}

	resp, err := h.backofficeSvc.List(c.Request().Context(), salesreturnparam.ListSalesReturnsRequest{
		Filter:     filter,
		Pagination: pagination(c),
	})
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) Get(c echo.Context) error {
	resp, err := h.backofficeSvc.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) Create(c echo.Context) error {
	const op = "salesreturnhandler.Create"

	var req salesreturnparam.CreateSalesReturnRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	if err := h.validator.ValidateCreateRequest(req); err != nil {
		return err
	}

	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.backofficeSvc.Create(c.Request().Context(), req, cl.EntityID)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) MarkRefunded(c echo.Context) error {
	const op = "salesreturnhandler.MarkRefunded"

	var req salesreturnparam.MarkRefundedRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	req.ID = c.Param("id")
	if err := h.validator.ValidateMarkRefundedRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeSvc.MarkRefunded(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

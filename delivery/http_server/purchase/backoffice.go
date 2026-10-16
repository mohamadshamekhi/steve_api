package purchasehandler

import (
	"strconv"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	purchaseparam "github.com/abolfazlnorzad/steve/param/purchase"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func pagination(c echo.Context) commonparam.PaginationRequest {
	pageNumber, err := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if err != nil {
		pageNumber = commonparam.DefaultPageNumber
	}
	pageSize, err := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if err != nil {
		pageSize = commonparam.DefaultPageSize
	}
	return commonparam.PaginationRequest{PageNumber: uint(pageNumber), PageSize: uint(pageSize)}
}

// parseDate reads an RFC3339 (or plain YYYY-MM-DD) query param. An unparsable
// value is treated as absent rather than as an error, so a malformed filter
// widens the report instead of breaking it.
func parseDate(raw string) *time.Time {
	if raw == "" {
		return nil
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return &t
	}
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		return &t
	}
	return nil
}

func (h Handler) List(c echo.Context) error {
	req := purchaseparam.ListPurchasesRequest{
		SupplierID:   c.QueryParam("supplier_id"),
		Status:       entity.PurchaseStatus(c.QueryParam("status")),
		PaymentState: entity.PaymentState(c.QueryParam("payment_state")),
		OrderID:      c.QueryParam("order_id"),
		From:         parseDate(c.QueryParam("from")),
		To:           parseDate(c.QueryParam("to")),
		Pagination:   pagination(c),
	}

	resp, err := h.backofficeSvc.List(c.Request().Context(), req)
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
	const op = "purchasehandler.Create"

	var req purchaseparam.CreatePurchaseRequest
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

func (h Handler) Update(c echo.Context) error {
	const op = "purchasehandler.Update"

	var req purchaseparam.UpdatePurchaseRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	req.ID = c.Param("id")

	if err := h.validator.ValidateUpdateRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeSvc.Update(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) Confirm(c echo.Context) error {
	resp, err := h.backofficeSvc.Confirm(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) Cancel(c echo.Context) error {
	resp, err := h.backofficeSvc.Cancel(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) SetPaidAmount(c echo.Context) error {
	const op = "purchasehandler.SetPaidAmount"

	var req purchaseparam.SetPaidAmountRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	req.ID = c.Param("id")

	resp, err := h.backofficeSvc.SetPaidAmount(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

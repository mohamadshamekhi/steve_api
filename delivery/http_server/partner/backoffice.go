package partnerhandler

import (
	"net/http"
	"strconv"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	partnerparam "github.com/abolfazlnorzad/steve/param/partner"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) success(c echo.Context, data any) error {
	lang := getLanguageFromHeader(c.Request().Header)
	msg := h.translator.Translate(lang, "global.action_successful", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(msg, data))
}

func (h Handler) List(c echo.Context) error {
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}
	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}
	req := partnerparam.ListRequest{
		Pagination: commonparam.PaginationRequest{PageSize: uint(pageSize), PageNumber: uint(pageNumber)},
		Status:     c.QueryParam("status"),
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

func (h Handler) Approve(c echo.Context) error {
	const op = "partnerhandler.Approve"
	var req partnerparam.ApproveRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	if err := h.validator.ValidateApproveRequest(req); err != nil {
		return err
	}
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.backofficeSvc.Approve(c.Request().Context(), c.Param("id"), cl.EntityID, req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) Reject(c echo.Context) error {
	const op = "partnerhandler.Reject"
	var req partnerparam.RejectRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	if err := h.validator.ValidateRejectRequest(req); err != nil {
		return err
	}
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.backofficeSvc.Reject(c.Request().Context(), c.Param("id"), cl.EntityID, req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) Suspend(c echo.Context) error {
	resp, err := h.backofficeSvc.Suspend(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) Unsuspend(c echo.Context) error {
	resp, err := h.backofficeSvc.Unsuspend(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) UpdatePricing(c echo.Context) error {
	const op = "partnerhandler.UpdatePricing"
	var req partnerparam.UpdatePricingRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	if err := h.validator.ValidateUpdatePricingRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeSvc.UpdatePricing(c.Request().Context(), c.Param("id"), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) SetProductConfig(c echo.Context) error {
	const op = "partnerhandler.SetProductConfig"
	var req partnerparam.SetProductConfigRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	resp, err := h.backofficeSvc.SetProductConfig(c.Request().Context(), c.Param("id"), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) BulkSetEligible(c echo.Context) error {
	const op = "partnerhandler.BulkSetEligible"
	var req partnerparam.BulkSetEligibleRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	resp, err := h.backofficeSvc.BulkSetEligible(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) RecordInPersonPayment(c echo.Context) error {
	const op = "partnerhandler.RecordInPersonPayment"
	var req partnerparam.RecordInPersonPaymentRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	resp, err := h.backofficeSvc.RecordInPersonPayment(c.Request().Context(), c.Param("id"), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) RecordHandover(c echo.Context) error {
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.backofficeSvc.RecordHandover(c.Request().Context(), c.Param("id"), cl.EntityID)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) ListProducts(c echo.Context) error {
	req := partnerparam.ProductListRequest{
		Pagination: pageFromQuery(c),
		Name:       c.QueryParam("name"),
		BrandID:    c.QueryParam("brand_id"),
		CategoryID: c.QueryParam("category_id"),
	}
	if e := c.QueryParam("eligible"); e == "true" || e == "false" {
		v := e == "true"
		req.Eligible = &v
	}
	resp, err := h.backofficeSvc.ListProducts(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) ListOrders(c echo.Context) error {
	req := partnerparam.OrderListRequest{
		Pagination: pageFromQuery(c),
		Status:     c.QueryParam("status"),
	}
	resp, err := h.backofficeSvc.ListOrders(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

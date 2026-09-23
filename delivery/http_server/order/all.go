package orderhandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

// tehranTZ is fixed at +03:30; Iran has no DST since 2022.
var tehranTZ = time.FixedZone("IRST", 3*3600+30*60)

func (h Handler) GetAll(c echo.Context) error {
	const op = "orderhandler.GetAll"
	var req orderparam.GetAllOrderRequest

	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}
	req.Pagination.PageSize = uint(pageSize)
	req.Pagination.PageNumber = uint(pageNumber)

	if err := (&echo.DefaultBinder{}).BindQueryParams(c, &req.Filter); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	// Dates are whole Tehran days; To is inclusive, so the bound is the next
	// midnight.
	if req.Filter.From != nil && *req.Filter.From != "" {
		t, pErr := time.ParseInLocation("2006-01-02", *req.Filter.From, tehranTZ)
		if pErr != nil {
			return richerror.New(op).WithKind(richerror.KindInvalid).
				WithMessageKey("validation.invalid_input").WithWarpError(pErr)
		}
		req.Filter.FromTime = &t
	}
	if req.Filter.To != nil && *req.Filter.To != "" {
		t, pErr := time.ParseInLocation("2006-01-02", *req.Filter.To, tehranTZ)
		if pErr != nil {
			return richerror.New(op).WithKind(richerror.KindInvalid).
				WithMessageKey("validation.invalid_input").WithWarpError(pErr)
		}
		t = t.AddDate(0, 0, 1)
		req.Filter.ToTime = &t
	}

	resp, err := h.backofficeSvc.GetAll(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
func (h Handler) GetUserOrder(c echo.Context) error {
	const op = "orderhandler.GetUserOrder"
	var req orderparam.GetUserOrderRequest

	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}
	req.Pagination.PageSize = uint(pageSize)
	req.Pagination.PageNumber = uint(pageNumber)
	req.UserID = c.Get("user_id").(string)
	resp, err := h.storefrontSvc.GetUserOrders(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) GetUserOrderStatics(c echo.Context) error {
	const op = "orderhandler.GetUserOrderStatics"
	var req orderparam.GetUserOrderStaticsRequest

	req.UserID = c.Get("user_id").(string)
	resp, err := h.storefrontSvc.GetUserOrderStatics(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

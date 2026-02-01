package shippinghandler

import (
	"net/http"
	"strconv"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	settingparam "github.com/abolfazlnorzad/steve/param/setting"
	"github.com/labstack/echo/v4"
)

func (h Handler) All(c echo.Context) error {
	const op = "shippinghandler.All"
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}
	isActive := c.QueryParam("is_active")
	var isActiveBool *bool
	if isActive != "" {
		b, _ := strconv.ParseBool(isActive)
		isActiveBool = &b

	} else {
		isActiveBool = nil
	}
	req := settingparam.GetAllShippingRequest{
		Pagination: commonparam.PaginationRequest{
			PageSize:   uint(pageSize),
			PageNumber: uint(pageNumber),
		},
		Filter: settingparam.FilterShipping{
			IsActive: isActiveBool,
		},
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

func (h Handler) GetAllStorefront(c echo.Context) error {
	const op = "shippinghandler.GetAllStorefront"

	req := settingparam.GetActiveShippingRequest{}
	resp, err := h.storefrontSvc.GetActiveShipping(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

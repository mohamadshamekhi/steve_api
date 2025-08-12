package discounthandler

import (
	"net/http"
	"strconv"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	"github.com/labstack/echo/v4"
)

func (h Handler) All(c echo.Context) error {
	const op = "discounthandler.All"
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}

	req := discountparam.GetAllDiscountsParamRequest{
		Pagination: commonparam.PaginationRequest{
			PageSize:   uint(pageSize),
			PageNumber: uint(pageNumber),
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

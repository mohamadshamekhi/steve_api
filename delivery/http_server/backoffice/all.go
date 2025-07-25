package backofficehandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h Handler) GetAll(c echo.Context) error {
	const op = "backofficehandler.GetAll"
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}
	phoneNumber := c.QueryParam("phone_number")
	req := backofficeparam.AllBackofficeRequest{
		Pagination: commonparam.PaginationRequest{
			PageSize:   uint(pageSize),
			PageNumber: uint(pageNumber),
		},
		Filter: commonparam.FilterRequest{
			Filter: []map[string]any{{"phone_number": phoneNumber}},
		},
	}

	resp, err := h.svc.GetAll(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

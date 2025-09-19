package sellhandler

import (
	"net/http"
	"strconv"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	sellparam "github.com/abolfazlnorzad/steve/param/sell"
	"github.com/labstack/echo/v4"
)

func (h Handler) All(c echo.Context) error {
	const op = "sellhandler.All"
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}
	userID := c.QueryParam("user_id")
	status := c.QueryParam("status")
	req := sellparam.GetAllSellParamRequest{
		Pagination: commonparam.PaginationRequest{
			PageSize:   uint(pageSize),
			PageNumber: uint(pageNumber),
		},
		Filter: commonparam.FilterRequest{
			Filter: []map[string]any{{"user_id": userID}, {"status": status}},
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

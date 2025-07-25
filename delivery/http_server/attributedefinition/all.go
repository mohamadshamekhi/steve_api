package attributedefinitionhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	attributedefinitionparam "github.com/abolfazlnorzad/steve/param/attribute_definition"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h Handler) All(c echo.Context) error {
	const op = "attributedefinitionhandler.All"
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}
	name := c.QueryParam("key")
	req := attributedefinitionparam.AllAttributeDefinitionRequest{
		Pagination: commonparam.PaginationRequest{
			PageSize:   uint(pageSize),
			PageNumber: uint(pageNumber),
		},
		Filter: commonparam.FilterRequest{
			Filter: []map[string]any{{"name": name}},
		},
	}

	resp, err := h.svc.All(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)

}

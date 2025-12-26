package sellhandler

import (
	"net/http"
	"strconv"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	sellparam "github.com/abolfazlnorzad/steve/param/sell"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/labstack/echo/v4"
)

func (h Handler) MySell(c echo.Context) error {
	const op = "sellhandler.MySell"
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}
	cl := claim.GetClaimsFromEchoContext(c)

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}
	req := sellparam.GetAllMySellParamRequest{
		Pagination: commonparam.PaginationRequest{
			PageSize:   uint(pageSize),
			PageNumber: uint(pageNumber),
		},
		Filter: commonparam.FilterRequest{
			Filter: []map[string]any{{"user_id": cl.EntityID}},
		},
	}

	resp, err := h.storefrontSvc.GetMySell(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)

}

func (h Handler) Statics(c echo.Context) error {
	const op = "sellhandler.Statics"

	req := sellparam.GetUserSellStaticsRequest{UserID: c.Get("user_id").(string)}
	resp, err := h.storefrontSvc.GetUserSellStatics(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)

}

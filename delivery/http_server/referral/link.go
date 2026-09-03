package referralhandler

import (
	"net/http"
	"strconv"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/labstack/echo/v4"
)

func pageFromQuery(c echo.Context) commonparam.PaginationRequest {
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}
	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}
	return commonparam.PaginationRequest{PageSize: uint(pageSize), PageNumber: uint(pageNumber)}
}

func intQuery(c echo.Context, name string) int {
	v, err := strconv.Atoi(c.QueryParam(name))
	if err != nil {
		return 0
	}
	return v
}

func (h Handler) success(c echo.Context, data any) error {
	lang := getLanguageFromHeader(c.Request().Header)
	msg := h.translator.Translate(lang, "global.action_successful", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(msg, data))
}

// GetMyLink returns the caller's permanent invite code, shareable link and
// progress toward the next reward.
func (h Handler) GetMyLink(c echo.Context) error {
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.svc.GetMyLink(c.Request().Context(), cl.EntityID)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

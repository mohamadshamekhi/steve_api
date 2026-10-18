package dashboardhandler

import (
	"strconv"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
	dashboardparam "github.com/abolfazlnorzad/steve/param/dashboard"
	"github.com/labstack/echo/v4"
)

// AISummary is the read-only, aggregated-only sales/inventory snapshot for
// the ai_agent role: GET /backoffice/dashboard/ai-summary?from=&to=&page_number=&page_size=
func (h Handler) AISummary(c echo.Context) error {
	pageNumber, err := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if err != nil {
		pageNumber = commonparam.DefaultPageNumber
	}
	pageSize, err := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if err != nil {
		pageSize = commonparam.DefaultPageSize
	}

	resp, err := h.svc.AISummary(c.Request().Context(), dashboardparam.AISummaryRequest{
		FromRaw: c.QueryParam("from"),
		ToRaw:   c.QueryParam("to"),
		Pagination: commonparam.PaginationRequest{
			PageNumber: uint(pageNumber),
			PageSize:   uint(pageSize),
		},
	})
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

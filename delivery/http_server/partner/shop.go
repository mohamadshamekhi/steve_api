package partnerhandler

import (
	"strconv"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
	partnerparam "github.com/abolfazlnorzad/steve/param/partner"
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

func (h Handler) ShopMe(c echo.Context) error {
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.shopSvc.Me(c.Request().Context(), cl.EntityID)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) CatalogList(c echo.Context) error {
	cl := claim.GetClaimsFromEchoContext(c)

	req := partnerparam.CatalogListRequest{Pagination: pageFromQuery(c)}
	if name := c.QueryParam("name"); name != "" {
		req.Filter.Filter = append(req.Filter.Filter, map[string]any{"name": name})
	}
	if brandID := c.QueryParam("brand_id"); brandID != "" {
		req.Filter.Filter = append(req.Filter.Filter, map[string]any{"brand_id": brandID})
	}
	if categoryID := c.QueryParam("category_id"); categoryID != "" {
		req.Filter.Filter = append(req.Filter.Filter, map[string]any{"category_id": categoryID})
	}

	resp, err := h.shopSvc.Catalog(c.Request().Context(), cl.EntityID, req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) MyOrders(c echo.Context) error {
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.shopSvc.MyOrders(c.Request().Context(), cl.EntityID, pageFromQuery(c))
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

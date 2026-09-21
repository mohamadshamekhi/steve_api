package inventoryhandler

import (
	"strconv"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/export"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

// defaultLowStockThreshold is used when ?low_stock=true comes without an
// explicit threshold.
const defaultLowStockThreshold = 3

func pagination(c echo.Context) commonparam.PaginationRequest {
	pageNumber, err := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if err != nil {
		pageNumber = commonparam.DefaultPageNumber
	}
	pageSize, err := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if err != nil {
		pageSize = commonparam.DefaultPageSize
	}
	return commonparam.PaginationRequest{PageNumber: uint(pageNumber), PageSize: uint(pageSize)}
}

// parseDate accepts RFC3339 or plain YYYY-MM-DD; an unparsable value is
// treated as absent so a malformed filter widens the report, never breaks it.
func parseDate(raw string) *time.Time {
	if raw == "" {
		return nil
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return &t
	}
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		return &t
	}
	return nil
}

func (h Handler) ListMovements(c echo.Context) error {
	req := inventoryparam.ListMovementsRequest{
		Filter: inventoryparam.MovementFilter{
			ProductID:  c.QueryParam("product_id"),
			VariantID:  c.QueryParam("variant_id"),
			SupplierID: c.QueryParam("supplier_id"),
			Type:       entity.StockMovementType(c.QueryParam("type")),
			RefType:    c.QueryParam("ref_type"),
			RefID:      c.QueryParam("ref_id"),
			From:       parseDate(c.QueryParam("from")),
			To:         parseDate(c.QueryParam("to")),
		},
		Pagination: pagination(c),
	}

	resp, err := h.backofficeSvc.ListMovements(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

// exportPageSize is the largest page the pagination helper accepts; a CSV
// export asks for it so the whole warehouse lands in one file.
const exportPageSize = 200

func (h Handler) InventoryValue(c echo.Context) error {
	req := inventoryparam.InventoryValueRequest{Pagination: pagination(c)}

	if c.QueryParam("low_stock") == "true" {
		req.LowStockOnly = true
		req.LowStockThreshold = defaultLowStockThreshold
		if raw := c.QueryParam("threshold"); raw != "" {
			if parsed, err := strconv.ParseUint(raw, 10, 64); err == nil {
				req.LowStockThreshold = parsed
			}
		}
	}

	// A CSV export covers the whole warehouse rather than the current page:
	// an انبارگردانی sheet split across pages is useless.
	if c.QueryParam("format") == "csv" {
		req.Pagination = commonparam.PaginationRequest{PageNumber: 1, PageSize: exportPageSize}
	}

	resp, err := h.backofficeSvc.InventoryValue(c.Request().Context(), req)
	if err != nil {
		return err
	}

	if c.QueryParam("format") == "csv" {
		rows := make([][]string, 0, len(resp.Rows)+1)
		for _, r := range resp.Rows {
			rows = append(rows, []string{r.ProductName, r.VariantID, string(r.StockMode),
				export.Uint(r.Quantity), export.Uint(r.AvgCost), export.Uint(r.Price),
				export.Uint(r.StockValue), export.Int(r.PotentialProfit)})
		}
		rows = append(rows, []string{"جمع کل", "", "", export.Uint(resp.TotalUnits), "", "",
			export.Uint(resp.TotalStockValue), ""})
		return export.CSV(c, "inventory-value",
			[]string{"کالا", "شناسه واریانت", "حالت موجودی", "موجودی", "قیمت خرید", "قیمت فروش", "ارزش انبار", "سود بالقوه"}, rows)
	}
	return h.success(c, resp)
}

func (h Handler) Adjust(c echo.Context) error {
	const op = "inventoryhandler.Adjust"

	var req inventoryparam.AdjustStockRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	if err := h.validator.ValidateAdjustRequest(req); err != nil {
		return err
	}

	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.backofficeSvc.AdjustStock(c.Request().Context(), req, cl.EntityID)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) SetStockConfig(c echo.Context) error {
	const op = "inventoryhandler.SetStockConfig"

	var req inventoryparam.SetStockConfigRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	if err := h.validator.ValidateSetStockConfigRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeSvc.SetStockConfig(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) ProfitReport(c echo.Context) error {
	req := inventoryparam.ProfitReportRequest{
		From:    parseDate(c.QueryParam("from")),
		To:      parseDate(c.QueryParam("to")),
		GroupBy: c.QueryParam("group_by"),
	}

	resp, err := h.backofficeSvc.ProfitReport(c.Request().Context(), req)
	if err != nil {
		return err
	}

	if c.QueryParam("format") == "csv" {
		rows := make([][]string, 0, len(resp.Rows)+1)
		for _, r := range resp.Rows {
			label := r.Label
			if label == "" {
				label = r.Key
			}
			rows = append(rows, []string{label, export.Int(r.UnitsSold),
				export.Uint(r.Revenue), export.Uint(r.Cost), export.Int(r.GrossProfit)})
		}
		rows = append(rows, []string{"جمع کل", export.Int(resp.Summary.UnitsSold),
			export.Uint(resp.Summary.Revenue), export.Uint(resp.Summary.Cost), export.Int(resp.Summary.GrossProfit)})
		return export.CSV(c, "profit-report",
			[]string{"عنوان", "تعداد فروش", "فروش", "بهای تمام‌شده", "سود ناخالص"}, rows)
	}
	return h.success(c, resp)
}

package producthandler

import (
	"fmt"
	"net/http"
	"time"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
)

const xlsxContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

// ExportExcel streams every product matching the given filters as an .xlsx file.
// Backoffice-only, gated behind the product:export permission (see routes.go).
// Accepts the same filter query params as the list endpoint (h.All), e.g.
// ?is_published=true to export only published products.
func (h Handler) ExportExcel(c echo.Context) error {
	const op = "producthandler.ExportExcel"

	ENTitle := c.QueryParam("en_title")
	name := c.QueryParam("name")
	isPub := c.QueryParam("is_published")
	isSecondHand := c.QueryParam("is_second_hand")
	brandID := c.QueryParam("brand_id")
	categoryID := c.QueryParam("category_id")
	req := productparam.ExportProductsRequest{
		Filter: commonparam.FilterRequest{
			Filter: []map[string]any{{"en_title": ENTitle}, {"name": name}, {"brand_id": brandID},
				{"category_id": categoryID}, {"is_published": isPub}, {"is_second_hand": isSecondHand}},
		},
	}

	resp, err := h.backofficeProductSvc.ExportProducts(c.Request().Context(), req)
	if err != nil {
		return err
	}

	file, err := buildProductsExcelFile(resp.Rows)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer file.Close()

	fileName := fmt.Sprintf("products-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Response().Header().Set(echo.HeaderContentType, xlsxContentType)
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Response().WriteHeader(http.StatusOK)

	return file.Write(c.Response())
}

func buildProductsExcelFile(rows []productparam.ProductExportRow) (*excelize.File, error) {
	const sheet = "Products"

	f := excelize.NewFile()
	if err := f.SetSheetName(f.GetSheetName(0), sheet); err != nil {
		return nil, err
	}

	headers := []string{
		"ID", "Name", "EN Title", "Slug", "Brand", "Category",
		"Is Published", "Is Second Hand", "Min Price", "Max Price",
		"Total Quantity", "Variant Count", "Sales Count", "Rating", "Created At",
	}
	for col, header := range headers {
		cell, err := excelize.CoordinatesToCellName(col+1, 1)
		if err != nil {
			return nil, err
		}
		if err := f.SetCellValue(sheet, cell, header); err != nil {
			return nil, err
		}
	}

	for i, row := range rows {
		r := i + 2 // header is row 1
		values := []any{
			row.ID, row.Name, row.EnTitle, row.Slug, row.Brand, row.Category,
			row.IsPublished, row.IsSecondHand, row.MinPrice, row.MaxPrice,
			row.TotalQuantity, row.VariantCount, row.SalesCount, row.Rating,
			row.CreatedAt.Format(time.RFC3339),
		}
		for col, value := range values {
			cell, err := excelize.CoordinatesToCellName(col+1, r)
			if err != nil {
				return nil, err
			}
			if err := f.SetCellValue(sheet, cell, value); err != nil {
				return nil, err
			}
		}
	}

	sheetIndex, err := f.GetSheetIndex(sheet)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(sheetIndex)

	return f, nil
}

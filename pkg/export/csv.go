// Package export renders report rows as a spreadsheet download.
package export

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// utf8BOM makes Excel open a UTF-8 file as UTF-8. Without it, Persian text
// in a CSV opens as mojibake on a default Windows install, which is how most
// of these reports get read.
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// CSV streams header + rows as an attachment named fileName (without
// extension), stamped with the current date.
func CSV(c echo.Context, fileName string, header []string, rows [][]string) error {
	name := fileName + "-" + time.Now().Format("2006-01-02") + ".csv"

	c.Response().Header().Set(echo.HeaderContentType, "text/csv; charset=utf-8")
	c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="`+name+`"`)
	c.Response().WriteHeader(http.StatusOK)

	if _, err := c.Response().Write(utf8BOM); err != nil {
		return err
	}

	w := csv.NewWriter(c.Response())
	if err := w.Write(header); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

// Uint formats an unsigned amount for a spreadsheet cell.
func Uint(v uint64) string { return strconv.FormatUint(v, 10) }

// Int formats a signed amount (a balance or a profit that can go negative).
func Int(v int64) string { return strconv.FormatInt(v, 10) }

// Float formats a percentage with two decimals.
func Float(v float64) string { return strconv.FormatFloat(v, 'f', 2, 64) }

// Date formats a timestamp as a plain day, which is what these reports are
// read by; a zero time renders empty rather than as year one.
func Date(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// DateTime keeps the clock time for ledger-style rows.
func DateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}

// Bool renders a flag the way a Persian-speaking reader expects.
func Bool(v bool) string {
	if v {
		return "بله"
	}
	return "خیر"
}

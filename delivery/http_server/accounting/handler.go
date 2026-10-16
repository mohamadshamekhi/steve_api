package accountinghandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	accountingservice "github.com/abolfazlnorzad/steve/service/accounting"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	accountingvalidator "github.com/abolfazlnorzad/steve/validator/accounting"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	backofficeSvc accountingservice.BackofficeService
	translator    *i18n.Translator
	validator     accountingvalidator.Validator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(
	backofficeSvc accountingservice.BackofficeService,
	translator *i18n.Translator,
	validator accountingvalidator.Validator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service,
) Handler {
	return Handler{
		backofficeSvc: backofficeSvc,
		translator:    translator,
		validator:     validator,
		authSvc:       authSvc,
		authCfg:       authCfg,
		aclSvc:        aclSvc,
	}
}

func getLanguageFromHeader(header http.Header) string {
	if header.Get("Accept-Language") == "en" {
		return "en"
	}
	return "fa"
}

func (h Handler) success(c echo.Context, data any) error {
	lang := getLanguageFromHeader(c.Request().Header)
	msg := h.translator.Translate(lang, "global.action_successful", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(msg, data))
}

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

// parseDate reads an RFC3339 (or plain YYYY-MM-DD) query param. An unparsable
// value is treated as absent rather than as an error, so a malformed filter
// widens the report instead of breaking it.
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

// parseDateEnd is parseDate for the upper bound of a window: a plain
// YYYY-MM-DD means "through the end of that day", not its first instant.
func parseDateEnd(raw string) *time.Time {
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		end := t.Add(24*time.Hour - time.Nanosecond)
		return &end
	}
	return parseDate(raw)
}

// parseBool reads an optional tri-state filter (set / unset / absent).
func parseBool(raw string) *bool {
	switch raw {
	case "true", "1":
		v := true
		return &v
	case "false", "0":
		v := false
		return &v
	}
	return nil
}

// wantsCSV reports whether the caller asked for a spreadsheet instead of JSON.
func wantsCSV(c echo.Context) bool {
	return c.QueryParam("format") == "csv"
}

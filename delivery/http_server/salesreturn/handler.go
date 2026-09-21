package salesreturnhandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	salesreturnservice "github.com/abolfazlnorzad/steve/service/salesreturn"
	salesreturnvalidator "github.com/abolfazlnorzad/steve/validator/salesreturn"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	backofficeSvc salesreturnservice.BackofficeService
	translator    *i18n.Translator
	validator     salesreturnvalidator.Validator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(
	backofficeSvc salesreturnservice.BackofficeService,
	translator *i18n.Translator,
	validator salesreturnvalidator.Validator,
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

// parseDate reads an RFC3339 (or plain YYYY-MM-DD) query param; an unparsable
// value widens the filter instead of failing the request.
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

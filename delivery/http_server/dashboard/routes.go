package dashboardhandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	ba := func() echo.MiddlewareFunc { return backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg) }

	g := e.Group("/backoffice/dashboard")
	// No single permission guards /overview: every section is filtered by
	// its own permission inside the service and is null when missing.
	g.GET("/overview", h.Overview, ba())
	// ai-summary is a single-purpose, hard-gated, aggregated-only endpoint —
	// see dashboardparam.AISummaryResponse for what it must never include.
	g.GET("/ai-summary", h.AISummary, ba(), backofficemiddleware.ACL(h.aclSvc, key.AIAgentReport))
}

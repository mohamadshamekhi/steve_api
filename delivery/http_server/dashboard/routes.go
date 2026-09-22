package dashboardhandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/backoffice/dashboard")
	// No single permission guards the route: every section is filtered by
	// its own permission inside the service and is null when missing.
	g.GET("/overview", h.Overview, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg))
}

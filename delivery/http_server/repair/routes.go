package repairhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/repair")
	g.POST("/storefront/add", h.Create, middleware.Auth(h.authSvc, h.authCfg))
	g.GET("/storefront/my-repair", h.All, middleware.Auth(h.authSvc, h.authCfg))
}

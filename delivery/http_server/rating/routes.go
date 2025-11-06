package ratinghandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/rating")
	g.POST("/storefront/upsert", h.Upsert, middleware.Auth(h.authSvc, h.authCfg))
}

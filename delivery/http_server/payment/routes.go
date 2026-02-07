package paymenthandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/payment")
	g.POST("/storefront/start", h.StartPay, middleware.Auth(h.authSvc, h.authCfg))
	g.POST("/storefront/verify", h.Verify, middleware.Auth(h.authSvc, h.authCfg))

}

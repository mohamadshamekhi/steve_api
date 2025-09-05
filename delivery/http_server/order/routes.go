package orderhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/order")
	g.POST("/register", h.Register, middleware.Auth(h.authSvc, h.authCfg))
	g.GET("/all", h.GetAll, middleware.Auth(h.authSvc, h.authCfg))

}

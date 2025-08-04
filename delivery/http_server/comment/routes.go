package commenthandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/comment")
	g.POST("/add", h.AddComment, middleware.Auth(h.authSvc, h.authCfg))
}

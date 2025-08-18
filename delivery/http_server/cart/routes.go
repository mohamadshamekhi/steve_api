package carthandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/cart")
	g.POST("/add-item", h.AddItemToCart, middleware.Auth(h.authSvc, h.authCfg))
	g.POST("/remove-item", h.RemoveItem, middleware.Auth(h.authSvc, h.authCfg))

}

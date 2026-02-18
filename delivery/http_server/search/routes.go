package searchhandler

import (
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/search")
	g.GET("/storefront", h.GetSearchStorefrontWithProductsAndArticles)
}

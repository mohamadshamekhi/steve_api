package contactushandler

import (
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/contact-us")
	g.POST("/storefront", h.Add)

}

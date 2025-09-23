package mediahandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/media")
	g.POST("/upload", h.UploadFile, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg))
	g.POST("/storefront/upload", h.StorefrontUploadFile, middleware.Auth(h.authSvc, h.authCfg))
	g.GET("/get-url/:id", h.GetURL)
}

package sellhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/sell")
	g.POST("/storefront/add", h.Create, middleware.Auth(h.authSvc, h.authCfg))
	g.GET("/storefront/my-sell", h.MySell, middleware.Auth(h.authSvc, h.authCfg))
	g.GET("/storefront/user-sell/statics", h.Statics, middleware.Auth(h.authSvc, h.authCfg))

	g.GET("/backoffice/:id", h.Get, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/backoffice/:id", h.Update, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.GET("/backoffice", h.All, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListProducts))
}

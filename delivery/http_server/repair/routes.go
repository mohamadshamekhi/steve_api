package repairhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/repair")
	g.POST("/storefront/add", h.Create, middleware.Auth(h.authSvc, h.authCfg))
	g.GET("/storefront/my-repair", h.MyRepair, middleware.Auth(h.authSvc, h.authCfg))

	g.GET("/backoffice/:id", h.Get, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/backoffice/:id", h.Update, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.GET("/backoffice", h.All, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListProducts))
}

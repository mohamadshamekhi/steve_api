package discounthandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/discount")
	g.POST("/create", h.Create, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.CreateDiscount))
	//g.GET("/:id", h.Get, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListDiscounts))
	g.GET("", h.All, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListDiscounts))
	g.PATCH("/:id", h.Update, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateDiscount))
}

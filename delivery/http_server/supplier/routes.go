package supplierhandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	ba := func() echo.MiddlewareFunc { return backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg) }

	g := e.Group("/backoffice/supplier")
	g.GET("", h.List, ba(), backofficemiddleware.ACL(h.aclSvc, key.ListSuppliers))
	g.GET("/:id", h.Get, ba(), backofficemiddleware.ACL(h.aclSvc, key.ViewSupplier))
	g.POST("", h.Create, ba(), backofficemiddleware.ACL(h.aclSvc, key.CreateSupplier))
	g.PATCH("/:id", h.Update, ba(), backofficemiddleware.ACL(h.aclSvc, key.UpdateSupplier))
	g.DELETE("/:id", h.Delete, ba(), backofficemiddleware.ACL(h.aclSvc, key.DeleteSupplier))
}

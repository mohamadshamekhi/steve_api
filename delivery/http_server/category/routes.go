package categoryhandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/category")
	g.POST("/create", h.Create, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.CreateCategory))
	g.GET("/:id", h.Get, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListCategories))
	g.GET("", h.All, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListCategories))
	g.PATCH("/:id", h.Update, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateCategory))
	g.DELETE("/:id", h.Delete, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.DeleteCategory))
	g.GET("/storefront", h.GetAllStorefront)
	g.GET("/storefront/get-filter/:category-id", h.GetFilter)
}

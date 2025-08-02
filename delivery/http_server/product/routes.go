package producthandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/product")
	g.POST("/create", h.Create, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.CreateProduct))
	g.GET("/:id", h.Get, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ViewProduct))
	g.GET("", h.All, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListProducts))
	g.PATCH("/:id", h.Update, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.DELETE("/:id", h.Delete, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.DeleteProduct))

	g.PATCH("/add-product-gallery", h.AddGallery, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/set-primary", h.SetPrimary, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/delete-product-gallery", h.DeleteGallery, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))

	g.PATCH("/:id/update-attrs", h.UpdateAttrs, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/:id/update-single-attr/:key", h.UpdateSingleAttrs, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))

	g.PATCH("/:id/add-variant", h.AddVariant, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/:id/update-variant/:varid", h.UpdateVariant, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/:id/delete-variant/:varid", h.DeleteVariant, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
}

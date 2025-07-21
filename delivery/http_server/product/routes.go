package producthandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/product")
	g.POST("/create", h.Create, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.CreateProduct))
	g.GET("/:id", h.Get, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ViewProduct))
	g.GET("", h.All, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListProducts))
	g.PATCH("/:id", h.Update, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.DELETE("/:id", h.Delete, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.DeleteProduct))

	g.PATCH("/add-product-gallery", h.AddGallery, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/set-primary", h.SetPrimary, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/delete-product-gallery", h.DeleteGallery, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))

	g.PATCH("/:id/update-attrs", h.UpdateAttrs, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/:id/update-single-attr/:key", h.UpdateSingleAttrs, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))

	g.PATCH("/:id/add-variant", h.AddVariant, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/:id/update-variant/:varid", h.UpdateVariant, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/:id/delete-variant/:varid", h.DeleteVariant, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
}

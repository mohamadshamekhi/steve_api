package producthandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/product")
	g.POST("/create", h.Create, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.CreateProduct))
	// Competitor-price agent. Registered before "/:id" for readability; Echo
	// matches static segments first either way.
	g.GET("/competitor-links", h.GetCompetitorLinks, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ManageCompetitorPrice))
	g.PATCH("/:id/competitor-price", h.UpdateCompetitorPrice, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ManageCompetitorPrice))
	g.GET("/:id", h.Get, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ViewProduct))
	g.GET("", h.All, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListProducts))
	g.GET("/export-excel", h.ExportExcel, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ExportProducts))
	g.PATCH("/:id", h.Update, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/publish/:id", h.Publish, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/un-publish/:id", h.UnPublish, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.DELETE("/:id", h.Delete, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.DeleteProduct))

	g.PATCH("/add-product-gallery", h.AddGallery, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/set-primary", h.SetPrimary, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/delete-product-gallery", h.DeleteGallery, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/set-gallery-alt", h.SetGalleryAlt, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))

	g.PATCH("/:id/update-attrs", h.UpdateAttrs, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/:id/update-single-attr/:key", h.UpdateSingleAttrs, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))

	g.PATCH("/:id/add-variant", h.AddVariant, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/:id/update-variant/:varid", h.UpdateVariant, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/:id/delete-variant/:varid", h.DeleteVariant, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))
	g.PATCH("/bulk-update-variants", h.BulkUpdateVariants, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateProduct))

	g.GET("/storefront/with-auto-discount", h.StorefrontGetProductsWithAutoDiscount)
	g.GET("/storefront", h.StorefrontGetFilterableProduct)
	g.GET("/storefront/get-mobile-list-wth-main-info", h.GetMobileList)
	g.GET("/storefront/get-products-with-special-discount", h.GetProductsWithSpecialDiscount)
	g.GET("/storefront/:id", h.StorefrontGetProduct)
	g.GET("/storefront/get-by-slug/:slug", h.StorefrontGetProductBySlug)
	// Public: target of the printed variant QR label; the variant ID is the
	// stable code (survives renames and slug changes).
	g.GET("/variant/:variant_id/resolve", h.ResolveVariant)
}

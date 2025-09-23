package articlehandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/article")
	g.POST("/create", h.Create, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.CreateArticle))
	g.GET("/:id", h.Get, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ViewArticle))
	g.GET("", h.All, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListArticles))
	g.PATCH("/:id", h.Update, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateArticle))
	g.DELETE("/:id", h.Delete, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.DeleteArticle))

	storefrontGroup := e.Group("article/storefront")
	storefrontGroup.GET("", h.GetAllStorefront)
	storefrontGroup.GET("/:id", h.GetSingleStorefront)

}

package articlehandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/article")
	g.POST("/create", h.Create, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.CreateArticle))
	g.GET("/:id", h.Get, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ViewArticle))
	g.GET("", h.All, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListArticles))
	g.PATCH("/:id", h.Update, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateArticle))
	g.DELETE("/:id", h.Delete, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.DeleteArticle))

	storefrontGroup := e.Group("article/storefront")
	storefrontGroup.GET("", h.GetAllStorefront)
	storefrontGroup.GET("/:id", h.GetSingleStorefront)

}

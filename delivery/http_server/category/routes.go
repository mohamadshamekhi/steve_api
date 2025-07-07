package categoryhandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/category")
	g.POST("/create", h.Create, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.CreateCategory))
	g.GET("/:id", h.Get)
	g.GET("", h.All)
	g.PATCH("/:id", h.Update, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateCategory))
	g.DELETE("/:id", h.Delete, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.DeleteCategory))
}

package commenthandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/comment")
	g.POST("/add", h.AddComment, middleware.Auth(h.authSvc, h.authCfg))
	g.POST("/update/:id", h.UpdateComment, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateComment))

	g.GET("", h.GetAllComments, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListComment))

	g.GET("/storefront/:commentabletype/:commentableid", h.GetComments)
}

package backofficehandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/back-office")
	g.GET("/", h.GetAll, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListBackofficeUsers))
	g.POST("/add-new-backoffice-user", h.AddNewBackofficeUser, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.CreateBackofficeUser))
	g.PATCH("/update/:id", h.UpdateBackofficeUser, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateBackofficeUser))
	g.GET("/me", h.MyProfile, middleware.Auth(h.authSvc, h.authCfg))
}

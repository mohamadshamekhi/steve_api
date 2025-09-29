package aclhandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/acl")
	g.GET("/permissions", h.GetAllPermissions, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListPermissions))
	g.GET("/roles", h.GetAllRole, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListRole))
}

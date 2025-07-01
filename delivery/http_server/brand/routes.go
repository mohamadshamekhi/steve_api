package brandhandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/brand")
	g.POST("/create", h.Create, backofficemiddleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.CanManageBrands))
}

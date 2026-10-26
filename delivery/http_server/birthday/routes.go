package birthdayhandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	b := e.Group("/backoffice/birthday", backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg))
	b.GET("/settings", h.GetSettings, backofficemiddleware.ACL(h.aclSvc, key.ViewBirthdayDiscount))
	b.PUT("/settings", h.UpdateSettings, backofficemiddleware.ACL(h.aclSvc, key.ManageBirthdayDiscount))
	b.GET("/rewards", h.ListRewards, backofficemiddleware.ACL(h.aclSvc, key.ViewBirthdayDiscount))
	b.POST("/run", h.Run, backofficemiddleware.ACL(h.aclSvc, key.ManageBirthdayDiscount))
}

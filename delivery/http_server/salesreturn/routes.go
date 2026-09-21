package salesreturnhandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	ba := func() echo.MiddlewareFunc { return backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg) }

	g := e.Group("/backoffice/sales-return")
	g.GET("", h.List, ba(), backofficemiddleware.ACL(h.aclSvc, key.ListSalesReturns))
	g.GET("/:id", h.Get, ba(), backofficemiddleware.ACL(h.aclSvc, key.ViewSalesReturn))
	g.POST("", h.Create, ba(), backofficemiddleware.ACL(h.aclSvc, key.CreateSalesReturn))
	g.POST("/:id/refund", h.MarkRefunded, ba(), backofficemiddleware.ACL(h.aclSvc, key.RefundSalesReturn))
}

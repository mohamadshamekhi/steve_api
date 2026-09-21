package purchasehandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	ba := func() echo.MiddlewareFunc { return backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg) }

	g := e.Group("/backoffice/purchase")
	g.GET("", h.List, ba(), backofficemiddleware.ACL(h.aclSvc, key.ListPurchases))
	g.GET("/:id", h.Get, ba(), backofficemiddleware.ACL(h.aclSvc, key.ViewPurchase))
	g.POST("", h.Create, ba(), backofficemiddleware.ACL(h.aclSvc, key.CreatePurchase))
	g.PATCH("/:id", h.Update, ba(), backofficemiddleware.ACL(h.aclSvc, key.UpdatePurchase))
	g.POST("/:id/confirm", h.Confirm, ba(), backofficemiddleware.ACL(h.aclSvc, key.ConfirmPurchase))
	g.POST("/:id/cancel", h.Cancel, ba(), backofficemiddleware.ACL(h.aclSvc, key.CancelPurchase))
	g.POST("/:id/payment", h.SetPaidAmount, ba(), backofficemiddleware.ACL(h.aclSvc, key.UpdatePurchase))
}

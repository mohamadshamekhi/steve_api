package orderdeliveryhandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	e.GET("/backoffice/order/:id/delivery-qrcode", h.GetQRCode,
		backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ViewOrderDetails))

	// Public: the customer arrives here by scanning the package QR and may
	// not be logged in; the token plus an OTP to the owner's phone is the proof.
	g := e.Group("/order/delivery")
	g.GET("/:token", h.GetInfo)
	g.POST("/:token/send-otp", h.SendOTP)
	g.POST("/:token/verify", h.Verify)
	g.GET("/:token/invoice", h.GetInvoice)
}

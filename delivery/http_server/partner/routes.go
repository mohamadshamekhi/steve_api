package partnerhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	ba := func() echo.MiddlewareFunc { return backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg) }

	// --- application entry (normal storefront user) ---
	p := e.Group("/partner")
	p.POST("/apply", h.Apply, middleware.Auth(h.authSvc, h.authCfg))
	p.GET("/me", h.Me, middleware.Auth(h.authSvc, h.authCfg))

	// --- partner shop panel (the partner's restricted BACKOFFICE account) ---
	shop := e.Group("/backoffice/partner-shop")
	shop.GET("/me", h.ShopMe, ba(), backofficemiddleware.ACL(h.aclSvc, key.PartnerShopAccess))
	shop.GET("/products", h.CatalogList, ba(), backofficemiddleware.ACL(h.aclSvc, key.PartnerShopAccess))
	shop.GET("/my-orders", h.MyOrders, ba(), backofficemiddleware.ACL(h.aclSvc, key.PartnerShopAccess))
	shop.GET("/cart", h.GetCart, ba(), backofficemiddleware.ACL(h.aclSvc, key.PartnerShopAccess))
	shop.POST("/cart/items", h.AddToCart, ba(), backofficemiddleware.ACL(h.aclSvc, key.PartnerShopOrder))
	shop.DELETE("/cart/items", h.RemoveFromCart, ba(), backofficemiddleware.ACL(h.aclSvc, key.PartnerShopOrder))
	shop.POST("/checkout", h.Checkout, ba(), backofficemiddleware.ACL(h.aclSvc, key.PartnerShopOrder))
	shop.POST("/payment/verify", h.VerifyPayment, ba(), backofficemiddleware.ACL(h.aclSvc, key.PartnerShopOrder))

	// --- partner management (admin staff: partner_manager) ---
	b := e.Group("/backoffice/partner")
	b.GET("", h.List, ba(), backofficemiddleware.ACL(h.aclSvc, key.ListPartners))
	b.GET("/:id", h.Get, ba(), backofficemiddleware.ACL(h.aclSvc, key.ViewPartner))
	b.POST("/:id/approve", h.Approve, ba(), backofficemiddleware.ACL(h.aclSvc, key.ApprovePartner))
	b.POST("/:id/reject", h.Reject, ba(), backofficemiddleware.ACL(h.aclSvc, key.RejectPartner))
	b.POST("/:id/suspend", h.Suspend, ba(), backofficemiddleware.ACL(h.aclSvc, key.SuspendPartner))
	b.POST("/:id/unsuspend", h.Unsuspend, ba(), backofficemiddleware.ACL(h.aclSvc, key.SuspendPartner))
	b.PATCH("/:id/pricing", h.UpdatePricing, ba(), backofficemiddleware.ACL(h.aclSvc, key.UpdatePartnerPricing))

	pp := e.Group("/backoffice/partner-product")
	pp.GET("", h.ListProducts, ba(), backofficemiddleware.ACL(h.aclSvc, key.ListPartners))
	pp.PATCH("/:id/config", h.SetProductConfig, ba(), backofficemiddleware.ACL(h.aclSvc, key.UpdatePartnerPricing))
	pp.POST("/bulk-eligible", h.BulkSetEligible, ba(), backofficemiddleware.ACL(h.aclSvc, key.UpdatePartnerPricing))

	// partner order fulfilment (partner_manager)
	po := e.Group("/backoffice/partner-order")
	po.GET("", h.ListOrders, ba(), backofficemiddleware.ACL(h.aclSvc, key.ListPartners))
	po.POST("/:id/in-person-payment", h.RecordInPersonPayment, ba(), backofficemiddleware.ACL(h.aclSvc, key.RecordPartnerHandover))
	po.POST("/:id/handover", h.RecordHandover, ba(), backofficemiddleware.ACL(h.aclSvc, key.RecordPartnerHandover))
}

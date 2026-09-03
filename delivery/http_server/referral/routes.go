package referralhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	// --- storefront: the member's own invite screen ---
	g := e.Group("/referral", middleware.Auth(h.authSvc, h.authCfg))
	g.GET("/me/link", h.GetMyLink)
	g.POST("/invite", h.InviteByPhone)
	g.GET("/me/referrals", h.ListMyReferrals)
	g.GET("/me/rewards", h.ListMyRewards)

	// --- backoffice: admin oversight of the referral program ---
	b := e.Group("/backoffice/referral", backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg))
	b.GET("/stats", h.AdminStats, backofficemiddleware.ACL(h.aclSvc, key.ListReferrals))
	b.GET("/settings", h.AdminGetSettings, backofficemiddleware.ACL(h.aclSvc, key.ListReferrals))
	b.PUT("/settings", h.AdminUpdateSettings, backofficemiddleware.ACL(h.aclSvc, key.ManageReferral))
	b.GET("", h.AdminList, backofficemiddleware.ACL(h.aclSvc, key.ListReferrals))
	b.GET("/:userID", h.AdminGet, backofficemiddleware.ACL(h.aclSvc, key.ViewReferral))
}

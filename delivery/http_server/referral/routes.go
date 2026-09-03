package referralhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/referral", middleware.Auth(h.authSvc, h.authCfg))

	g.GET("/me/link", h.GetMyLink)
	g.POST("/invite", h.InviteByPhone)
	g.GET("/me/referrals", h.ListMyReferrals)
	g.GET("/me/rewards", h.ListMyRewards)
}

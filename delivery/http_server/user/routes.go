package userhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/user")

	g.POST("/send-otp", h.Login)
	g.POST("/verify-otp", h.VerifyOTPCode)
	g.PATCH("/:id/update-profile", h.UpdateProfile, middleware.Auth(h.authSvc, h.authCfg))
	g.POST("/add-address", h.AddAddress, middleware.Auth(h.authSvc, h.authCfg))
	g.POST("/update-address/:id", h.UpdateAddress, middleware.Auth(h.authSvc, h.authCfg))

	g.GET("", h.GetAllUser, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.acl, key.ListUsers))
}

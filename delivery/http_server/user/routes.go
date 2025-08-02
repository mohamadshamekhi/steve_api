package userhandler

import (
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/user")

	g.POST("/send-otp", h.Login)
	g.POST("/verify-otp", h.VerifyOTPCode)
	g.PATCH("/:id/update-profile", h.UpdateProfile)
}

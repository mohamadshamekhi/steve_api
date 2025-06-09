package backofficeauthhandler

import "github.com/labstack/echo/v4"

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/back-office")
	g.POST("/send-otp", h.LoginByPhoneNumber)
	g.POST("/verify-otp", h.VerifyOTPCode)
}

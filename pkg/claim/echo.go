package claim

import (
	"github.com/abolfazlnorzad/steve/config"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	"github.com/labstack/echo/v4"
)

func GetClaimsFromEchoContext(c echo.Context) *authservice.Claims {
	return c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)
}

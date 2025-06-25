package backofficemiddleware

import (
	"github.com/abolfazlnorzad/steve/config"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	mw "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Auth(service authservice.Service, cfg authservice.Config) echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{
		ContextKey: config.AuthMiddlewareContextKey,
		SigningKey: []byte(cfg.SignKey),
		// TODO  - as sign method string to config
		SigningMethod: "HS256",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims, err := service.ParseToken(auth)
			if err != nil {
				return nil, err
			}

			return claims, nil
		},
	})
}

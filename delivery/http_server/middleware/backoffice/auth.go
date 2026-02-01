package backofficemiddleware

import (
	"context"

	"github.com/abolfazlnorzad/steve/config"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	mw "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func BackofficeAuth(service authservice.Service, cfg authservice.Config) echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{
		ContextKey: config.AuthMiddlewareContextKey,
		ErrorHandler: func(c echo.Context, err error) error {
			return richerror.New("middleware.auth").WithMessageKey("auth.you_must_logged_in").WithKind(richerror.KindUnauthorized)
		},
		SigningKey: []byte(cfg.SignKey),
		// TODO  - as sign method string to config
		SigningMethod: "HS256",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims, err := service.ParseToken(auth)
			if err != nil {
				return nil, err
			}
			if claims.Kind != "backoffice" {
				return nil, richerror.New("backofficemiddleware.BackofficeAuth").WithMessageKey("auth.you_must_logged_in").WithKind(richerror.KindForbidden)
			}
			res, err := service.CheckTokenIsInBlockList(c.Request().Context(), claims.Kind, claims.ID)
			if err != nil {
				return nil, err
			}
			if !res {
				ctx := context.WithValue(
					c.Request().Context(),
					"backoffice_id",
					claims.EntityID,
				)
				c.SetRequest(c.Request().WithContext(ctx))
				return claims, nil
			}
			return nil, richerror.New("middleware.auth").WithMessageKey("auth.you_must_logged_in").WithKind(richerror.KindUnauthorized)
		},
	})
}

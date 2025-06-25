package backofficemiddleware

import (
	"github.com/abolfazlnorzad/steve/pkg/claim"
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
	"github.com/labstack/echo/v4"
)

func IsActive(service backofficeservice.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := claim.GetClaimsFromEchoContext(c)
			resp, err := service.CheckBackOfficeUserIsActive(c.Request().Context(), claims.EntityID)
			if err != nil || !resp {
				return err
			}
			return next(c)
		}
	}
}

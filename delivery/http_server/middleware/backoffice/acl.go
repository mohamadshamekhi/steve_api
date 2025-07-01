package backofficemiddleware

import (
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	"github.com/labstack/echo/v4"
)

func ACL(service aclservice.Service,
	permissions ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			claims := claim.GetClaimsFromEchoContext(c)
			isAllowed, err := service.CheckAccess(c.Request().Context(), claims.EntityID, claims.RoleKeys, claims.PermissionKeys, permissions...)
			if err != nil {
				// TODO - log unexpected error
				return err
			}

			if !isAllowed {
				return richerror.New("backofficemiddleware.ACL").WithMessageKey("global.forbidden").WithKind(richerror.KindForbidden)
			}

			return next(c)
		}
	}
}

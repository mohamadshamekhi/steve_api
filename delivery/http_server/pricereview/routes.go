package pricereviewhandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	auth := backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg)
	viewAcl := backofficemiddleware.ACL(h.aclSvc, key.ViewPriceReview)
	manageAcl := backofficemiddleware.ACL(h.aclSvc, key.ManagePriceReview)

	settingGroup := e.Group("/backoffice/setting/price-review")
	settingGroup.GET("", h.GetSetting, auth, viewAcl)
	settingGroup.PUT("", h.UpdateSetting, auth, manageAcl)

	reviewGroup := e.Group("/backoffice/price-review")
	reviewGroup.GET("/list", h.List, auth, viewAcl)
	reviewGroup.POST("/:id/confirm", h.Confirm, auth, manageAcl)
}

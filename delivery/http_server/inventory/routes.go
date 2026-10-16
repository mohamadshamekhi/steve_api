package inventoryhandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	ba := func() echo.MiddlewareFunc { return backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg) }

	g := e.Group("/backoffice/inventory")
	g.GET("/movements", h.ListMovements, ba(), backofficemiddleware.ACL(h.aclSvc, key.ViewInventory))
	g.GET("/stock", h.InventoryValue, ba(), backofficemiddleware.ACL(h.aclSvc, key.ViewInventoryReport))
	g.POST("/adjust", h.Adjust, ba(), backofficemiddleware.ACL(h.aclSvc, key.AdjustInventory))
	g.POST("/stocktake", h.Stocktake, ba(), backofficemiddleware.ACL(h.aclSvc, key.AdjustInventory))
	g.POST("/backfill-sales", h.BackfillSales, ba(), backofficemiddleware.ACL(h.aclSvc, key.AdjustInventory))
	g.PATCH("/stock-config", h.SetStockConfig, ba(), backofficemiddleware.ACL(h.aclSvc, key.UpdateStockConfig))

	r := e.Group("/backoffice/report")
	r.GET("/profit", h.ProfitReport, ba(), backofficemiddleware.ACL(h.aclSvc, key.ViewProfitReport))
}

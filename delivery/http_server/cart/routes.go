package carthandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/cart")
	g.POST("/add-item", h.AddItemToCart, middleware.Auth(h.authSvc, h.authCfg))
	g.POST("/remove-item", h.RemoveItem, middleware.Auth(h.authSvc, h.authCfg))
	g.GET("/get-active", h.GetActive, middleware.Auth(h.authSvc, h.authCfg))
	g.POST("/calc-cart", h.CalculateCart, middleware.Auth(h.authSvc, h.authCfg))
	g.POST("/delete-item", h.DeleteItem, middleware.Auth(h.authSvc, h.authCfg))
	g.POST("/accept-price-changing", h.AcceptPriceChanging, middleware.Auth(h.authSvc, h.authCfg))

	g.GET("/backoffice/all", h.GetAll, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListCarts))
	g.GET("/backoffice/cart/:id", h.GetSingle, middleware.Auth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ViewCartDetails))

}

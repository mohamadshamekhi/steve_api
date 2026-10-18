package backofficehandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	g := e.Group("/back-office")
	g.GET("", h.GetAll, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ListBackofficeUsers))
	g.POST("/add-new-backoffice-user", h.AddNewBackofficeUser, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.CreateBackofficeUser))
	g.PATCH("/update/:id", h.UpdateBackofficeUser, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.UpdateBackofficeUser))
	g.GET("/me", h.MyProfile, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg))
	g.POST("/logout", h.Logout, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg))
	g.GET("/:id", h.Get, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.ViewBackofficeUser))
	g.POST("/:id/issue-service-token", h.IssueServiceToken, backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg), backofficemiddleware.ACL(h.aclSvc, key.IssueServiceToken))
}

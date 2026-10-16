package accountinghandler

import (
	backofficemiddleware "github.com/abolfazlnorzad/steve/delivery/http_server/middleware/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	ba := func() echo.MiddlewareFunc { return backofficemiddleware.BackofficeAuth(h.authSvc, h.authCfg) }

	// هزینه‌های عملیاتی
	expense := e.Group("/backoffice/expense")
	expense.GET("", h.ListExpenses, ba(), backofficemiddleware.ACL(h.aclSvc, key.ListExpenses))
	expense.GET("/:id", h.GetExpense, ba(), backofficemiddleware.ACL(h.aclSvc, key.ViewExpense))
	expense.POST("", h.CreateExpense, ba(), backofficemiddleware.ACL(h.aclSvc, key.CreateExpense))
	expense.PATCH("/:id", h.UpdateExpense, ba(), backofficemiddleware.ACL(h.aclSvc, key.UpdateExpense))
	expense.DELETE("/:id", h.DeleteExpense, ba(), backofficemiddleware.ACL(h.aclSvc, key.DeleteExpense))

	// صندوق و نقدینگی
	cash := e.Group("/backoffice/cash")
	cash.GET("", h.ListCashTransactions, ba(), backofficemiddleware.ACL(h.aclSvc, key.ViewCash))
	cash.POST("", h.CreateCashTransaction, ba(), backofficemiddleware.ACL(h.aclSvc, key.CreateCash))
	cash.DELETE("/:id", h.DeleteCashTransaction, ba(), backofficemiddleware.ACL(h.aclSvc, key.DeleteCash))

	// گزارش‌های حسابداری
	report := e.Group("/backoffice/report")
	report.GET("/cash", h.CashReport, ba(), backofficemiddleware.ACL(h.aclSvc, key.ViewCashReport))
	report.GET("/profit-loss", h.ProfitLoss, ba(), backofficemiddleware.ACL(h.aclSvc, key.ViewProfitLossReport))
	report.GET("/supplier-ledger", h.SupplierLedger, ba(), backofficemiddleware.ACL(h.aclSvc, key.ViewSupplierLedger))
}

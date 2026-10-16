package dashboardhandler

import (
	dashboardparam "github.com/abolfazlnorzad/steve/param/dashboard"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/labstack/echo/v4"
)

// Overview is the backoffice home page: GET /backoffice/dashboard/overview?from=YYYY-MM-DD&to=YYYY-MM-DD
func (h Handler) Overview(c echo.Context) error {
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.svc.Overview(c.Request().Context(), dashboardparam.OverviewRequest{
		FromRaw:          c.QueryParam("from"),
		ToRaw:            c.QueryParam("to"),
		BackofficeUserID: cl.EntityID,
	})
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

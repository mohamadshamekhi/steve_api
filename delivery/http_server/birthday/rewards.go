package birthdayhandler

import (
	birthdayparam "github.com/abolfazlnorzad/steve/param/birthday"
	"github.com/labstack/echo/v4"
)

// ListRewards returns the granted birthday codes, newest first.
// GET /backoffice/birthday/rewards?page_number&page_size
func (h Handler) ListRewards(c echo.Context) error {
	resp, err := h.svc.ListRewards(c.Request().Context(), birthdayparam.ListRewardsRequest{Pagination: pageFromQuery(c)})
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

// Run triggers the birthday job now instead of waiting for the daily schedule.
// Idempotent: users already rewarded for this birthday are skipped.
// POST /backoffice/birthday/run
func (h Handler) Run(c echo.Context) error {
	resp, err := h.svc.Run(c.Request().Context())
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

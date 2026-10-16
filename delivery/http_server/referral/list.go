package referralhandler

import (
	referralparam "github.com/abolfazlnorzad/steve/param/referral"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/labstack/echo/v4"
)

// ListMyReferrals returns the paginated list of people who registered through
// the caller's invite link (plus outstanding SMS invites).
func (h Handler) ListMyReferrals(c echo.Context) error {
	cl := claim.GetClaimsFromEchoContext(c)
	req := referralparam.ListMyReferralsRequest{Pagination: pageFromQuery(c)}
	resp, err := h.svc.ListMyReferrals(c.Request().Context(), cl.EntityID, req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

// ListMyRewards returns the discount codes the caller has earned so far.
func (h Handler) ListMyRewards(c echo.Context) error {
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.svc.ListMyRewards(c.Request().Context(), cl.EntityID)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

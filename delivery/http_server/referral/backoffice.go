package referralhandler

import (
	referralparam "github.com/abolfazlnorzad/steve/param/referral"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

// AdminList returns the referrer leaderboard (sorted by total joins, desc).
// GET /backoffice/referral?page_number&page_size&min_joined
func (h Handler) AdminList(c echo.Context) error {
	req := referralparam.AdminListReferralsRequest{
		Pagination: pageFromQuery(c),
		MinJoined:  intQuery(c, "min_joined"),
	}
	resp, err := h.backofficeSvc.List(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

// AdminGet returns one referrer's account, paginated invitees and reward codes.
// GET /backoffice/referral/:userID?page_number&page_size
func (h Handler) AdminGet(c echo.Context) error {
	req := referralparam.AdminGetReferralRequest{
		UserID:     c.Param("userID"),
		Pagination: pageFromQuery(c),
	}
	resp, err := h.backofficeSvc.Get(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

// AdminStats returns the program-wide dashboard summary.
// GET /backoffice/referral/stats
func (h Handler) AdminStats(c echo.Context) error {
	resp, err := h.backofficeSvc.Stats(c.Request().Context())
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

// AdminGetSettings returns the current effective reward/welcome configuration.
// GET /backoffice/referral/settings
func (h Handler) AdminGetSettings(c echo.Context) error {
	resp, err := h.backofficeSvc.GetSettings(c.Request().Context())
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

// AdminUpdateSettings replaces the reward/welcome configuration.
// PUT /backoffice/referral/settings
func (h Handler) AdminUpdateSettings(c echo.Context) error {
	const op = "referralhandler.AdminUpdateSettings"
	var req referralparam.UpdateReferralSettingsRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	if err := h.validator.ValidateUpdateReferralSettingsRequest(req); err != nil {
		return err
	}
	req.AdminUserID = claim.GetClaimsFromEchoContext(c).EntityID

	resp, err := h.backofficeSvc.UpdateSettings(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

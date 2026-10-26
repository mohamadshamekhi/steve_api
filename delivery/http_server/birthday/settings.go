package birthdayhandler

import (
	birthdayparam "github.com/abolfazlnorzad/steve/param/birthday"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

// GetSettings returns the effective birthday-discount configuration.
// GET /backoffice/birthday/settings
func (h Handler) GetSettings(c echo.Context) error {
	resp, err := h.svc.GetSettings(c.Request().Context())
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

// UpdateSettings replaces the birthday-discount configuration.
// PUT /backoffice/birthday/settings
func (h Handler) UpdateSettings(c echo.Context) error {
	const op = "birthdayhandler.UpdateSettings"
	var req birthdayparam.UpdateSettingsRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	if err := h.validator.ValidateUpdateSettingsRequest(req); err != nil {
		return err
	}
	req.AdminUserID = claim.GetClaimsFromEchoContext(c).EntityID

	resp, err := h.svc.UpdateSettings(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

package backofficehandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) Logout(c echo.Context) error {
	const op = "backofficehandler.Logout"
	var req backofficeparam.LogoutRequest
	cliam := claim.GetClaimsFromEchoContext(c)
	req.Kind = cliam.Kind
	req.ExpireAt = cliam.ExpiresAt.Time
	req.JWTID = cliam.ID
	resp, err := h.svc.Logout(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

package userhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/labstack/echo/v4"
)

func (h Handler) Logout(c echo.Context) error {
	const op = "userhandler.Logout"
	var req userparam.LogoutRequest
	cliam := claim.GetClaimsFromEchoContext(c)
	req.Kind = cliam.Kind
	req.ExpireAt = cliam.ExpiresAt.Time
	req.JWTID = cliam.ID
	resp, err := h.storeFrontSvc.Logout(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

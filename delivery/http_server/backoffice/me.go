package backofficehandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	claim2 "github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) MyProfile(c echo.Context) error {
	const op = "backofficehandler.myProfile"
	var req backofficeparam.MeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	claim := claim2.GetClaimsFromEchoContext(c)
	req.ID = claim.EntityID
	if err := h.vld.ValidateMeRequest(req); err != nil {
		return err
	}
	resp, err := h.svc.GetMyProfile(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

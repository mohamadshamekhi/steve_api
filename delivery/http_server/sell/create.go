package sellhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	sellparam "github.com/abolfazlnorzad/steve/param/sell"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) Create(c echo.Context) error {
	const op = "sellhandler.Create"
	var req sellparam.AddSellRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	cl := claim.GetClaimsFromEchoContext(c)
	req.UserID = cl.EntityID

	if err := h.validator.ValidateAddSellRequest(req); err != nil {
		return err
	}
	resp, err := h.storefrontSvc.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

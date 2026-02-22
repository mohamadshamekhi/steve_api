package carthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetActive(c echo.Context) error {
	const op = "carthandler.GetActive"
	cl := claim.GetClaimsFromEchoContext(c)
	req := cartparam.MyActiveCartRequest{
		UserID: cl.EntityID,
	}

	resp, err := h.storefrontSvc.GetMyActiveCart(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) GetActiveByUserIDBackoffice(c echo.Context) error {
	const op = "carthandler.GetActiveByUserIDBackoffice"
	var req cartparam.MyActiveCartBackofficeRequest
	req.UserID = c.Param("id")

	resp, err := h.backofficeSvc.GetMyActiveCartBackoffice(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

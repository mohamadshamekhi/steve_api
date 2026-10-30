package producthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetCompetitorLinks(c echo.Context) error {
	resp, err := h.backofficeProductSvc.GetCompetitorLinks(c.Request().Context(), productparam.GetCompetitorLinksRequest{})
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(successMessage, resp))
}

func (h Handler) UpdateCompetitorPrice(c echo.Context) error {
	const op = "producthandler.UpdateCompetitorPrice"

	var req productparam.UpdateCompetitorPriceRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	req.ProductID = c.Param("id")

	if err := h.vld.ValidateUpdateCompetitorPriceRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeProductSvc.UpdateCompetitorPrice(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(successMessage, resp))
}

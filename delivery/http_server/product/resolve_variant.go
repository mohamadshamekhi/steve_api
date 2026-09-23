package producthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/labstack/echo/v4"
)

func (h Handler) ResolveVariant(c echo.Context) error {
	resp, err := h.storefrontProductSvc.ResolveVariant(c.Request().Context(), productparam.ResolveVariantRequest{VariantID: c.Param("variant_id")})
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(h.translator.Translate(lang, "global.action_successful", nil), resp))
}

package producthandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) Get(c echo.Context) error {
	const op = "producthandler.Get"
	id := c.Param("id")
	req := productparam.GetProductRequest{
		ID: id,
	}

	if err := h.vld.ValidateGetProductRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeProductSvc.GetProduct(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

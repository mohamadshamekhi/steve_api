package carthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) AddItemToCart(c echo.Context) error {
	const op = "carthandler.AddItemToCart"
	var req cartparam.AddItemToCartRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	cl := claim.GetClaimsFromEchoContext(c)
	c.Set("user_id", cl.EntityID)

	if err := h.vld.ValidateAddItemToCartRequest(req); err != nil {
		return err
	}
	resp, err := h.storefrontSvc.AddItemToCart(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

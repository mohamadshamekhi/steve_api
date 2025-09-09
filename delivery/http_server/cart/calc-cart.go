package carthandler

import (
	"context"
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) CalculateCart(c echo.Context) error {
	const op = "handler.cart.CalculateCart"

	claims := claim.GetClaimsFromEchoContext(c)

	userID := claims.EntityID

	var req cartparam.CalculateCartRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("global.invalid_input").
			WithWarpError(err)
	}

	if err := h.vld.ValidateCalcCartRequest(req); err != nil {
		return err
	}
	ctx := context.WithValue(c.Request().Context(), "user_id", userID)
	activeCart, err := h.storefrontSvc.GetMyActiveCart(ctx, cartparam.MyActiveCartRequest{UserID: userID})
	if err != nil {
		return err
	}
	var cartItems []entity.CartItem
	for _, item := range activeCart.Cart.Items {
		cartItems = append(cartItems, entity.CartItem{
			ProductID:         item.ProductID,
			VariantID:         item.VariantID,
			Quantity:          item.Quantity,
			ProductName:       item.ProductName,
			VariantAttributes: item.VariantAttributes,
			PriceAtTime:       item.PriceAtTime,
		})
	}
	resp, err := h.orderSvc.FinalizeCartForCheckout(ctx, entity.Cart{
		ID:     activeCart.Cart.ID,
		UserID: activeCart.Cart.User.ID,
		Items:  cartItems,
		Status: entity.CartStatus(activeCart.Cart.Status),
	}, req.DiscountCode)
	if err != nil {
		return err
	}

	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

package partnerhandler

import (
	partnerparam "github.com/abolfazlnorzad/steve/param/partner"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetCart(c echo.Context) error {
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.shopSvc.GetCart(c.Request().Context(), cl.EntityID)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) AddToCart(c echo.Context) error {
	const op = "partnerhandler.AddToCart"
	var req partnerparam.AddToCartRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.shopSvc.AddToCart(c.Request().Context(), cl.EntityID, req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) RemoveFromCart(c echo.Context) error {
	const op = "partnerhandler.RemoveFromCart"
	var req partnerparam.RemoveFromCartRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.shopSvc.RemoveFromCart(c.Request().Context(), cl.EntityID, req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) Checkout(c echo.Context) error {
	const op = "partnerhandler.Checkout"
	var req partnerparam.CheckoutRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.shopSvc.Checkout(c.Request().Context(), cl.EntityID, req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) VerifyPayment(c echo.Context) error {
	const op = "partnerhandler.VerifyPayment"
	var req partnerparam.VerifyPaymentRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.shopSvc.VerifyPayment(c.Request().Context(), cl.EntityID, req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

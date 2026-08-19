package producthandler

import (
	"errors"
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) AddVariant(c echo.Context) error {
	const op = "productHandler.AddVariant"
	var req productparam.AddVariantRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	req.ProductID = c.Param("id")

	if err := h.vld.ValidateAddVariantRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeProductSvc.AddVariant(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) DeleteVariant(c echo.Context) error {
	const op = "productHandler.DeleteVariant"
	var req productparam.DeleteVariantRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	req.ProductID = c.Param("id")
	req.VariantID = c.Param("varid")

	if err := h.vld.ValidateDeleteVariantRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeProductSvc.DeleteVariant(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

// BulkUpdateVariants updates a batch of variants (possibly across different
// products) in one request. Each item goes through the exact same
// validate-then-update path as the single UpdateVariant endpoint - it's just
// looped - so one bad item (unknown variant, bad field, ...) is reported in
// that item's result instead of failing the whole batch.
func (h Handler) BulkUpdateVariants(c echo.Context) error {
	const op = "productHandler.BulkUpdateVariants"
	var req productparam.BulkUpdateVariantsRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	if len(req.Items) == 0 {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.required").
			WithTemplateData(map[string]any{"Attribute": "items"})
	}

	lang := getLanguageFromHeader(c.Request().Header)
	ctx := c.Request().Context()

	results := make([]productparam.BulkUpdateVariantItemResult, 0, len(req.Items))
	for _, item := range req.Items {
		if vErr := h.vld.ValidateUpdateVariantRequest(item); vErr != nil {
			results = append(results, productparam.BulkUpdateVariantItemResult{
				VariantID: item.VariantID,
				Success:   false,
				Message:   h.translateError(lang, vErr),
			})
			continue
		}

		if _, err := h.backofficeProductSvc.UpdateVariant(ctx, item); err != nil {
			results = append(results, productparam.BulkUpdateVariantItemResult{
				VariantID: item.VariantID,
				Success:   false,
				Message:   h.translateError(lang, err),
			})
			continue
		}

		results = append(results, productparam.BulkUpdateVariantItemResult{
			VariantID: item.VariantID,
			Success:   true,
		})
	}

	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, productparam.BulkUpdateVariantsResponse{Results: results})

	return c.JSON(http.StatusOK, apiResponse)
}

// translateError turns any error into a single display string for a bulk
// item result. richerror carries the i18n message key; anything else (a
// plain error we didn't expect) falls back to the generic message key.
func (h Handler) translateError(lang string, err error) string {
	var re *richerror.RichError
	if errors.As(err, &re) {
		return h.translator.Translate(lang, re.MessageKey(), re.TemplateData())
	}
	return h.translator.Translate(lang, "global.something_wrong", nil)
}

func (h Handler) UpdateVariant(c echo.Context) error {
	const op = "productHandler.UpdateVariant"
	var req productparam.UpdateVariantRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	req.ProductID = c.Param("id")
	req.VariantID = c.Param("varid")

	if err := h.vld.ValidateUpdateVariantRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeProductSvc.UpdateVariant(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

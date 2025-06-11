package middleware

import (
	"errors"
	"fmt"
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewCentralErrorHandler(translator *i18n.Translator) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		lang := getLanguageFromHeader(c.Request().Header)
		apiResponse, code := parseError(err, translator, lang)
		logger.L().Error(fmt.Sprintf("request field : %v", err))
		if err := c.JSON(code, apiResponse); err != nil {
			logger.L().Error("failed to send json error response, error: %v\n", err)
		}
	}
}
func parseError(err error, translator *i18n.Translator, lang string) (response.APIResponse, int) {
	var richErr *richerror.RichError
	if errors.As(err, &richErr) {
		code := mapKindToHTTPStatusCode(richErr.Kind())
		mainMessage := translator.Translate(lang, richErr.MessageKey(), richErr.TemplateData())

		if fieldErrorKeys, ok := richErr.GetMeta()["fields"].(map[string]string); ok {
			translatedErrors := make(map[string]string)

			for field, key := range fieldErrorKeys {
				attributeKey := "attr." + field
				translatedAttribute := translator.Translate(lang, attributeKey, nil)

				templateData := map[string]any{
					"Attribute": translatedAttribute,
				}
				finalMessage := translator.Translate(lang, key, templateData)

				translatedErrors[field] = finalMessage
			}

			return response.APIResponse{Message: mainMessage, Errors: translatedErrors}, code
		}

		if code >= http.StatusInternalServerError {
			msg := translator.Translate(lang, "global.something_wrong", nil)
			return response.APIResponse{Message: msg}, code
		}

		msg := translator.Translate(lang, richErr.MessageKey(), richErr.TemplateData())
		return response.APIResponse{Message: msg}, code
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		return response.APIResponse{Message: "درخواست شما با خطا مواجه شد."}, echoErr.Code
	}

	msg := translator.Translate(lang, "global.something_wrong", nil)
	return response.APIResponse{Message: msg}, http.StatusInternalServerError
}
func mapKindToHTTPStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindBadRequest:
		return http.StatusBadRequest
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindInvalid:
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}

func getLanguageFromHeader(header http.Header) string {
	if header.Get("Accept-Language") == "en" {
		return "en"
	}
	return "fa"
}

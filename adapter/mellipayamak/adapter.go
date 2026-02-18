package mellipayamak

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

const (
	PatternOTPKey     = "otp"
	SuccessPaymentKey = "payment"
	endpointPattern   = "SendByBaseNumber"
	endpointSimple    = "SendSimpleSMS2"
)

type Config struct {
	Username   string            `koanf:"username"`
	Password   string            `koanf:"password"`
	URL        string            `koanf:"url"`
	Patterns   map[string]string `koanf:"patterns"`
	FromNumber string            `koanf:"from_number"`
}

type Adapter struct {
	cfg    Config
	logger *slog.Logger
	client *http.Client
}

func New(cfg Config) Adapter {
	return Adapter{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type SendOTPResponse struct {
	Value map[string]string `json:"Value"`
}

func (a Adapter) SendOTP(ctx context.Context, to string, code string) error {
	const OP = "adapter.mellipayamak.SendOTP"

	patternID, ok := a.cfg.Patterns[PatternOTPKey]
	if !ok {
		return richerror.New(OP).
			WithWarpError(fmt.Errorf("pattern id for %s not found", PatternOTPKey)).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("")
	}
	body := map[string]any{
		"username": a.cfg.Username,
		"password": a.cfg.Password,
		"text":     []string{code},
		"to":       to,
		"bodyId":   patternID,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return richerror.New(OP).WithWarpError(err).WithMessageKey("failed to marshal request")
	}

	baseURL := strings.TrimSuffix(a.cfg.URL, "/")
	finalURL := fmt.Sprintf("%s/%s", baseURL, endpointPattern)

	req, err := http.NewRequestWithContext(ctx, "POST", finalURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return richerror.New(OP).WithWarpError(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return richerror.New(OP).WithWarpError(err)
	}

	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var apiResp SendOTPResponse
	if err := json.Unmarshal(bodyBytes, &apiResp.Value); err != nil {
		logger.L().Error("sms_non_json_response", "body", string(bodyBytes))
		return richerror.New(OP).WithMessageKey("invalid response format")
	}

	key := "d"
	if value, ok := apiResp.Value[key]; ok {
		if len(value) > 15 {
			logger.L().Info("sms_sent_successfully", "to", to, "response_id", apiResp.Value[key])
			return nil
		} else {
			logger.L().Error("sms_sent_unsuccessfully", "to", to, "response_id", apiResp.Value[key])
			return richerror.New(OP).
				WithMeta(map[string]any{"response": apiResp.Value}).
				WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
		}
	}
	logger.L().Error("sms_provider_logic_error",
		"response", apiResp.Value,
		"to", to,
	)
	return richerror.New(OP).
		WithMeta(map[string]any{"response": apiResp.Value}).
		WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)

}

func (a Adapter) SendPayment(ctx context.Context, to string, text []string) error {
	const OP = "adapter.mellipayamak.SendOTP"

	patternID, ok := a.cfg.Patterns[SuccessPaymentKey]
	if !ok {
		return richerror.New(OP).
			WithWarpError(fmt.Errorf("pattern id for %s not found", SuccessPaymentKey)).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("")
	}
	body := map[string]any{
		"username": a.cfg.Username,
		"password": a.cfg.Password,
		"text":     text,
		"to":       to,
		"bodyId":   patternID,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return richerror.New(OP).WithWarpError(err).WithMessageKey("failed to marshal request")
	}

	baseURL := strings.TrimSuffix(a.cfg.URL, "/")
	finalURL := fmt.Sprintf("%s/%s", baseURL, endpointPattern)

	req, err := http.NewRequestWithContext(ctx, "POST", finalURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return richerror.New(OP).WithWarpError(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return richerror.New(OP).WithWarpError(err)
	}

	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var apiResp SendOTPResponse
	if err := json.Unmarshal(bodyBytes, &apiResp.Value); err != nil {
		logger.L().Error("sms_non_json_response", "body", string(bodyBytes))
		return richerror.New(OP).WithMessageKey("invalid response format")
	}

	key := "d"
	if value, ok := apiResp.Value[key]; ok {
		if len(value) > 15 {
			logger.L().Info("sms_sent_successfully", "to", to, "response_id", apiResp.Value[key])
			return nil
		} else {
			logger.L().Error("sms_sent_unsuccessfully", "to", to, "response_id", apiResp.Value[key])
			return richerror.New(OP).
				WithMeta(map[string]any{"response": apiResp.Value}).
				WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
		}
	}
	logger.L().Error("sms_provider_logic_error",
		"response", apiResp.Value,
		"to", to,
	)
	return richerror.New(OP).
		WithMeta(map[string]any{"response": apiResp.Value}).
		WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)

}

type SendSimpleResponse struct {
	XMLName xml.Name `xml:"string"`
	Value   string   `xml:",chardata"`
}

func (a Adapter) SendSimple(ctx context.Context, to string, message string) error {
	const OP = "adapter.mellipayamak.SendSimple"

	baseURL := strings.TrimSuffix(a.cfg.URL, "/")
	finalURL := fmt.Sprintf("%s/%s", baseURL, endpointSimple)
	// endpointSimple = "post/Send.asmx/SendSimpleSMS2"

	form := url.Values{}
	form.Set("username", a.cfg.Username)
	form.Set("password", a.cfg.Password)
	form.Set("from", a.cfg.FromNumber)
	form.Set("to", to)
	form.Set("text", message)
	form.Set("isflash", "false")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, finalURL, strings.NewReader(form.Encode()))
	if err != nil {
		return richerror.New(OP).WithWarpError(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/xml")

	resp, err := a.client.Do(req)
	if err != nil {
		return richerror.New(OP).WithWarpError(err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var result SendSimpleResponse
	if err := xml.Unmarshal(bodyBytes, &result); err != nil {
		logger.L().Error("sms_non_xml_response", "body", string(bodyBytes))
		return richerror.New(OP).
			WithMessageKey("invalid response format")
	}

	codeStr := strings.TrimSpace(result.Value)

	code, _ := strconv.Atoi(codeStr)

	if code == 1 || len(codeStr) > 15 {
		logger.L().Info("sms_sent_successfully", "to", to, "response_id", code)
		return nil
	}

	logger.L().Error("sms_provider_logic_error", "to", to, "response_id", code)
	return richerror.New(OP).
		WithMessageKey("sms_provider_logic_error").
		WithMeta(map[string]any{"code": code}).
		WithKind(richerror.KindUnexpected)
}

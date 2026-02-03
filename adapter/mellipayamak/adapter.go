package mellipayamak

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

const (
	PatternOTPKey   = "otp"
	endpointPattern = "SendByBaseNumber"
	endpointSimple  = "SendSimpleSMS"
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
	return a.doRequest(ctx, OP, endpointPattern, body, to)
}

func (a Adapter) SendSimple(ctx context.Context, to string, message string) error {
	const OP = "adapter.mellipayamak.SendSimple"

	body := map[string]any{
		"username": a.cfg.Username,
		"password": a.cfg.Password,
		"from":     a.cfg.FromNumber,
		"to":       to,
		"text":     message,
		"isFlash":  false,
	}

	return a.doRequest(ctx, OP, endpointSimple, body, to)
}

type MelliPayamakResponse struct {
	Value        string `json:"Value"`
	RetStatus    int    `json:"RetStatus"`
	StrRetStatus string `json:"StrRetStatus"`
}

func (a Adapter) doRequest(ctx context.Context, op string, endpoint string, body any, to string) error {

	jsonData, err := json.Marshal(body)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("failed to marshal request")
	}

	baseURL := strings.TrimSuffix(a.cfg.URL, "/")
	finalURL := fmt.Sprintf("%s/%s", baseURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, "POST", finalURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return richerror.New(op).WithWarpError(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return richerror.New(op).WithWarpError(err)
	}

	defer func() { _ = resp.Body.Close() }()

	bodyBytes, _ := io.ReadAll(resp.Body)
	resultStr := strings.TrimSpace(string(bodyBytes))
	var apiResp MelliPayamakResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		logger.L().Error("sms_non_json_response", "body", string(bodyBytes))
		return richerror.New(op).WithMessageKey("invalid response format")
	}

	if apiResp.RetStatus != 1 && apiResp.RetStatus != 0 {

		logger.L().Error("sms_provider_logic_error",
			"status", apiResp.RetStatus,
			"message", apiResp.StrRetStatus,
			"to", to,
		)
		return richerror.New(op).
			WithMeta(map[string]any{"ret_status": apiResp.RetStatus, "msg": apiResp.StrRetStatus}).
			WithMessageKey("provider rejected data")
	}
	if apiResp.RetStatus == 1 {
		logger.L().Info("sms_sent_successfully", "to", to, "sms_id", resultStr)
		return nil
	}
	logger.L().Error("sms_sent_unsuccessfully", "to", to, "sms_id", resultStr)
	return richerror.New(op).
		WithMeta(map[string]any{"ret_status": apiResp.RetStatus, "msg": apiResp.StrRetStatus}).
		WithMessageKey("sms_sent_unsuccessfully")
}

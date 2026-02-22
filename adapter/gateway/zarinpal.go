package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type Config struct {
	MerchantID  string `koanf:"merchant_id"`
	VerifyURL   string `koanf:"verify_url"`
	StartPayURL string `koanf:"start_pay_url"`
	RequestURL  string `koanf:"request_url"`
	Currency    string `koanf:"currency"`
}

type ZarinpalAdapter struct {
	cfg    Config
	client *http.Client
}

func NewZarinpalAdapter(cfg Config) *ZarinpalAdapter {
	return &ZarinpalAdapter{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type requestPayload struct {
	MerchantID  string   `json:"merchant_id"`
	Amount      uint64   `json:"amount"` // ریال
	CallbackURL string   `json:"callback_url"`
	Description string   `json:"description"`
	Metadata    metadata `json:"metadata,omitempty"`
	Currency    string   `json:"currency"`
}

type metadata struct {
	Mobile string `json:"mobile,omitempty"`
	Email  string `json:"email,omitempty"`
}
type responsePayload struct {
	Data struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		Authority string `json:"authority"`
		Fee       int    `json:"fee"`
	} `json:"data"`
	Errors any `json:"errors"`
}

func (z *ZarinpalAdapter) RequestPayment(amount uint64, callbackURL string, description string) (string, string, error) {
	const op = "adapter.zarinpal.RequestPayment"
	payload := requestPayload{
		MerchantID:  z.cfg.MerchantID,
		Amount:      amount,
		CallbackURL: callbackURL,
		Description: "test",
		Currency:    z.cfg.Currency,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", "", richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	req, err := http.NewRequest("POST", z.cfg.RequestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := z.client.Do(req)
	if err != nil {
		return "", "", richerror.New(op).
			WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("gateway.connection_error")
	}
	defer resp.Body.Close()

	var result responsePayload
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("333", err)
		return "", "", richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.Data.Code != 100 {
		return "", "", richerror.New(op).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("gateway.request_failed").
			WithMeta(map[string]any{
				"gateway_code":    result.Data.Code,
				"gateway_message": result.Data.Message,
			})
	}

	paymentURL := fmt.Sprintf("%s/%s", z.cfg.StartPayURL, result.Data.Authority)

	return result.Data.Authority, paymentURL, nil
}

type verifyResponsePayload struct {
	Data struct {
		Code     int    `json:"code"`
		Message  string `json:"message"`
		CardHash string `json:"card_hash"`
		CardPan  string `json:"card_pan"`
		RefID    int64  `json:"ref_id"`
		Fee      int    `json:"fee"`
	} `json:"data"`
	Errors interface{} `json:"errors"`
}

func (z *ZarinpalAdapter) VerifyPayment(amount uint64, authority string) (string, string, error) {
	const op = "adapter.zarinpal.VerifyPayment"

	//amountInRials := amount * 10

	payload := map[string]interface{}{
		"merchant_id": z.cfg.MerchantID,
		"amount":      amount,
		"authority":   authority,
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", z.cfg.VerifyURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", richerror.New(op).WithWarpError(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.client.Do(req)
	if err != nil {
		return "", "", richerror.New(op).WithWarpError(err).WithMessageKey("gateway.connection_error")
	}
	defer resp.Body.Close()

	var result verifyResponsePayload
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", richerror.New(op).WithWarpError(err)
	}

	fmt.Println("result", result)
	if result.Data.Code == 100 || result.Data.Code == 101 {
		return fmt.Sprintf("%d", result.Data.RefID), result.Data.CardPan, nil
	}

	errCode := result.Data.Code
	errMsg := result.Data.Message

	return "", "", richerror.New(op).
		WithKind(richerror.KindInvalid).
		WithMessageKey("gateway.verification_failed").
		WithTemplateData(map[string]any{
			"gateway_code":    errCode,
			"gateway_message": errMsg,
		})
}

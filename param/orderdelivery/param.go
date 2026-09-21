package orderdeliveryparam

type GetQRCodeRequest struct {
	OrderID string `param:"id"`
}

type GetQRCodeResponse struct {
	OrderID string `json:"order_id"`
	URL     string `json:"url"`
	// QRCode is a data URI (image/png;base64) ready for an <img src>.
	QRCode string `json:"qr_code"`
}

type GetInfoRequest struct {
	Token string `param:"token"`
}

type GetInfoResponse struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	StatusLabel string `json:"status_label"`
	MaskedPhone string `json:"masked_phone"`
	CanConfirm  bool   `json:"can_confirm"`
}

type SendOTPRequest struct {
	Token       string `param:"token" json:"-"`
	PhoneNumber string `json:"phone_number"`
}

type SendOTPResponse struct {
	MaskedPhone string `json:"masked_phone"`
}

type VerifyRequest struct {
	Token       string `param:"token" json:"-"`
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

type VerifyResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

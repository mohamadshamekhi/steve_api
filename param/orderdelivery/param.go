package orderdeliveryparam

import "time"

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
	OrderID string  `json:"order_id"`
	Status  string  `json:"status"`
	Invoice Invoice `json:"invoice"`
	// InvoiceToken unlocks GET /order/delivery/:token/invoice (header
	// X-Invoice-Token) for InvoiceTokenExpiresIn seconds, so a page refresh
	// doesn't need a new OTP.
	InvoiceToken          string `json:"invoice_token"`
	InvoiceTokenExpiresIn int64  `json:"invoice_token_expires_in"`
}

type GetInvoiceRequest struct {
	Token        string
	InvoiceToken string
}

type GetInvoiceResponse struct {
	Invoice Invoice `json:"invoice"`
}

type Invoice struct {
	OrderID         string         `json:"order_id"`
	Status          string         `json:"status"`
	StatusLabel     string         `json:"status_label"`
	Items           []InvoiceItem  `json:"items"`
	Subtotal        uint           `json:"subtotal"`
	ShippingCost    uint           `json:"shipping_cost"`
	Discount        uint           `json:"discount"`
	FinalAmount     uint           `json:"final_amount"`
	ShippingAddress InvoiceAddress `json:"shipping_address"`
	CreatedAt       time.Time      `json:"created_at"`
}

type InvoiceItem struct {
	ProductName  string `json:"product_name"`
	VariantLabel string `json:"variant_label"`
	Quantity     uint   `json:"quantity"`
	UnitPrice    uint   `json:"unit_price"`
	Discount     uint   `json:"discount"`
	Total        uint   `json:"total"`
}

type InvoiceAddress struct {
	Title    string `json:"title"`
	Province string `json:"province"`
	City     string `json:"city"`
	Street   string `json:"street"`
	Postcode string `json:"postcode"`
}

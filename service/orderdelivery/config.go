package orderdeliveryservice

import "time"

type Config struct {
	// BaseURL is the storefront page the QR points to; the token is appended
	// as ?token=<token>.
	BaseURL        string        `koanf:"base_url"`
	OTPPrefix      string        `koanf:"otp_prefix"`
	MaxOTPAttempts int64         `koanf:"max_otp_attempts"`
	AttemptWindow  time.Duration `koanf:"attempt_window"`
	QRSize         int           `koanf:"qr_size"`
	// InvoiceTokenTTL is how long the post-OTP invoice token stays valid.
	InvoiceTokenTTL time.Duration `koanf:"invoice_token_ttl"`
}

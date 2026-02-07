package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusSuccessful PaymentStatus = "successful"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)

type Payment struct {
	ID                   string         `bson:"_id" json:"id"`
	UserID               string         `bson:"user_id" json:"user_id"`
	OrderID              string         `bson:"order_id" json:"order_id"`
	Amount               uint64         `bson:"amount" json:"amount"`
	Gateway              string         `bson:"gateway" json:"gateway"`
	Status               PaymentStatus  `bson:"status" json:"status"`
	GatewayTransactionID string         `bson:"gateway_transaction_id" json:"gateway_transaction_id"`
	TrackingCode         string         `bson:"tracking_code,omitempty" json:"tracking_code,omitempty"`
	MaskedCardNumber     string         `bson:"masked_card_number,omitempty" json:"masked_card_number,omitempty"`
	GatewayResponseCode  int            `bson:"gateway_response_code,omitempty" json:"-"`
	GatewayData          map[string]any `bson:"gateway_data,omitempty" json:"gateway_data,omitempty"`
	ErrorMessage         string         `bson:"error_message,omitempty" json:"error_message,omitempty"`
	Description          string         `bson:"description" json:"description"`
	CreatedAt            time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time      `bson:"updated_at" json:"updated_at"`
}

func (p *Payment) BeforeCreate() {
	p.ID = ulid.Make().String()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	p.Status = PaymentStatusPending
}

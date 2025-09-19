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
	ID string `bson:"_id"`

	OrderID string `bson:"order_id"`

	Amount uint `bson:"amount"`

	PaymentMethod string `bson:"payment_method"`

	Status PaymentStatus `bson:"status"`

	TransactionID string `bson:"transaction_id,omitempty"`
	Gateway       string `bson:"gateway,omitempty"`

	ProviderData map[string]interface{} `bson:"provider_data,omitempty"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (p *Payment) BeforeCreate() {
	p.ID = ulid.Make().String()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	p.Status = PaymentStatusPending
}

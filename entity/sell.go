package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type SellStatus string

const (
	// وضعیت‌های اولیه و در حال بررسی
	SellStatusNew             SellStatus = "new"              // ثبت شده (منتظر بررسی اولیه کارشناس)
	SellStatusPendingShipment SellStatus = "pending_shipment" // در انتظار ارسال (تایید اولیه، منتظر ارسال دستگاه توسط کاربر)
	SellStatusShipped         SellStatus = "shipped"          // ارسال شده (کاربر دستگاه را ارسال کرده)
	SellStatusUnderInspection SellStatus = "under_inspection" // در حال بررسی فنی (دستگاه دریافت شده و در حال کارشناسی است)
	SellStatusOfferMade       SellStatus = "offer_made"       // پیشنهاد قیمت داده شده (منتظر تایید یا رد کاربر)

	// وضعیت‌های نهایی موفق
	SellStatusOfferAccepted    SellStatus = "offer_accepted"    // پیشنهاد پذیرفته شد (کاربر قیمت را تایید کرده، منتظر پرداخت)
	SellStatusPaymentProcessed SellStatus = "payment_processed" // پرداخت انجام شد (فرآیند با موفقیت به پایان رسیده)

	// وضعیت‌های نهایی ناموفق یا لغو شده
	SellStatusOfferRejected      SellStatus = "offer_rejected"       // پیشنهاد رد شد (کاربر قیمت را رد کرده، دستگاه باید بازگردانده شود)
	SellStatusReturnedToCustomer SellStatus = "returned_to_customer" // به مشتری بازگردانده شد
	SellStatusCanceled           SellStatus = "canceled"             // لغو شده (توسط کاربر یا ادمین در مراحل اولیه)
	SellStatusRejected           SellStatus = "rejected"             // رد شده (درخواست اولیه توسط کارشناس رد شده است)
)

type Sell struct {
	ID                string     `bson:"_id"`
	UserID            string     `bson:"user_id"`
	BrandID           string     `bson:"brand_id"`
	ProductID         string     `bson:"product_id"`
	Color             string     `bson:"color"`
	BatteryHealth     string     `bson:"battery_health"`
	RegisterStatus    string     `bson:"register_status"`
	PhysicalCondition string     `bson:"physical_condition"`
	Description       string     `bson:"description"`
	Status            SellStatus `bson:"status"`
	CreatedAt         time.Time  `bson:"created_at"`
	UpdatedAt         time.Time  `bson:"updated_at"`
}

func (s *Sell) BeforeCreate() {
	s.ID = ulid.Make().String()
	s.CreatedAt = time.Now()
	s.UpdatedAt = s.CreatedAt
}

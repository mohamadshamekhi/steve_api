package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type OrderStatus string

const (
	// OrderStatusPendingPayment (در انتظار پرداخت)
	// سفارش ثبت شده اما هنوز پرداخت نشده است. موجودی انبار برای این سفارش رزرو شده است.
	OrderStatusPendingPayment OrderStatus = "pending_payment"

	// OrderStatusProcessing (در حال پردازش)
	// پرداخت با موفقیت انجام شده و سفارش در حال آماده‌سازی در انبار است.
	OrderStatusProcessing OrderStatus = "processing"

	// OrderStatusShipped (ارسال شده)
	// بسته به شرکت حمل و نقل (مانند پست) تحویل داده شده و کد پیگیری برای آن ثبت شده است.
	OrderStatusShipped OrderStatus = "shipped"

	// OrderStatusDelivered (تحویل داده شده)
	// بسته با موفقیت به دست مشتری رسیده است. این معمولاً آخرین وضعیت یک سفارش موفق است.
	OrderStatusDelivered OrderStatus = "delivered"

	// OrderStatusCancelled (لغو شده)
	// سفارش توسط کاربر یا ادمین لغو شده است. موجودی کالاهای این سفارش باید به انبار برگردانده شود.
	OrderStatusCancelled OrderStatus = "cancelled"

	// OrderStatusRefunded (مرجوع شده)
	// کالا توسط مشتری مرجوع شده و هزینه آن به صورت کامل یا جزئی بازگشت داده شده است.
	OrderStatusRefunded OrderStatus = "refunded"
)

type PurchaseType string

const (
	PurchaseTypeInPerson PurchaseType = "in_person" // حضوری
	PurchaseTypeOnline   PurchaseType = "online"    // آنلاین
)

type Order struct {
	ID               string            `bson:"_id"`
	UserID           string            `bson:"user_id"`
	CartID           string            `bson:"cart_id"`
	Items            []OrderItem       `bson:"items"`
	ShippingAddress  Address           `bson:"shipping_address"`
	SubtotalAmount   uint              `bson:"subtotal_amount"`
	ShippingAmount   uint              `bson:"shipping_amount"`
	AppliedDiscounts []AppliedDiscount `bson:"applied_discounts,omitempty"`
	FinalAmount      uint              `bson:"final_amount"`
	Shipping         Shipping          `bson:"shipping,omitempty"`
	Status           OrderStatus       `bson:"status"`
	PurchaseType     PurchaseType      `bson:"purchase_type"`
	CreatedAt        time.Time         `bson:"created_at"`
	UpdatedAt        time.Time         `bson:"updated_at"`
}

func (o *Order) BeforeCreate() {
	o.ID = ulid.Make().String()
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()

	if o.PurchaseType == "" {
		o.PurchaseType = PurchaseTypeOnline
	}
}

type OrderItem struct {
	ProductID         string             `bson:"product_id"`
	VariantID         string             `bson:"variant_id"`
	Quantity          uint               `bson:"quantity"`
	ProductName       string             `bson:"product_name"`
	VariantAttributes []VariantAttribute `bson:"variant_attributes"`

	PriceAtTime uint `bson:"price_at_time"`

	//DiscountAmount uint `bson:"discount_amount,omitempty"`
}

type AppliedDiscount struct {
	DiscountID   string `bson:"discount_id" json:"discount_id"`
	DiscountCode string `bson:"discount_code,omitempty" json:"discount_code,omitempty"`
	Amount       int64  `bson:"amount" json:"amount"`
}

func CanTransitionOrderStatus(from, to OrderStatus) bool {
	allowed := map[OrderStatus][]OrderStatus{
		OrderStatusProcessing: {
			OrderStatusShipped,
			OrderStatusCancelled,
			OrderStatusDelivered,
		},

		OrderStatusShipped: {
			OrderStatusDelivered,
			OrderStatusProcessing,
		},
		OrderStatusPendingPayment: {},
		OrderStatusDelivered:      {},
		OrderStatusCancelled:      {},
		OrderStatusRefunded:       {},
	}

	for _, st := range allowed[from] {
		if st == to {
			return true
		}
	}

	return false
}

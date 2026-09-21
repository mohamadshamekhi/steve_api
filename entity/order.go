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
	// ManualDiscountAmount is the discount applied to the invoice as a whole
	// at the counter, on top of the per-line discounts, in Toman.
	ManualDiscountAmount uint         `bson:"manual_discount_amount,omitempty"`
	FinalAmount          uint         `bson:"final_amount"`
	Shipping             Shipping     `bson:"shipping,omitempty"`
	Status               OrderStatus  `bson:"status"`
	PurchaseType         PurchaseType `bson:"purchase_type"`
	// Partner (colleague-sales) order metadata. Set when the order is placed
	// from the partner panel; used for financial reporting and the in-person
	// handover flow.
	IsPartnerOrder bool       `bson:"is_partner_order,omitempty"`
	PartnerID      string     `bson:"partner_id,omitempty"`
	HandoverBy     string     `bson:"handover_by,omitempty"`
	HandoverAt     *time.Time `bson:"handover_at,omitempty"`
	// DeliveryToken is the unguessable secret printed as a QR code on the
	// package; the customer scans it and confirms receipt with an OTP.
	// Minted lazily the first time the backoffice asks for the QR.
	DeliveryToken       string     `bson:"delivery_token,omitempty"`
	DeliveredAt         *time.Time `bson:"delivered_at,omitempty"`
	DeliveryConfirmedBy string     `bson:"delivery_confirmed_by,omitempty"`
	// ThankYouClaimed guards the one-time post-purchase thank-you coupon so a
	// replayed status change never mints a second one.
	ThankYouClaimed    bool   `bson:"thank_you_claimed,omitempty"`
	ThankYouDiscountID string `bson:"thank_you_discount_id,omitempty"`
	CreatedAt      time.Time  `bson:"created_at"`
	UpdatedAt      time.Time  `bson:"updated_at"`
}

func (o *Order) BeforeCreate() {
	o.ID = ulid.Make().String()
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()

	if o.PurchaseType == "" {
		o.PurchaseType = PurchaseTypeOnline
	}
}

// SupplySource records where the units of one order line actually came from.
type SupplySource string

const (
	// SupplySourceStock — از موجودی مغازه تأمین شد.
	SupplySourceStock SupplySource = "stock"
	// SupplySourceExternal — در لحظه از تأمین‌کننده/مغازه دیگر تهیه شد.
	SupplySourceExternal SupplySource = "external"
	// SupplySourceMixed — بخشی از انبار، بخشی تأمین بیرونی (backorder).
	SupplySourceMixed SupplySource = "mixed"
	// SupplySourceService — خدمات؛ تأمینی در کار نیست.
	SupplySourceService SupplySource = "service"
)

type OrderItem struct {
	ProductID         string             `bson:"product_id"`
	VariantID         string             `bson:"variant_id"`
	Quantity          uint               `bson:"quantity"`
	ProductName       string             `bson:"product_name"`
	VariantAttributes []VariantAttribute `bson:"variant_attributes"`

	PriceAtTime uint `bson:"price_at_time"`

	// DiscountAmount is everything taken off this line: campaign/coupon
	// discounts plus the manual discount given at the counter, in Toman for
	// the whole line.
	DiscountAmount uint `bson:"discount_amount,omitempty"`

	// ManualDiscount is the part of DiscountAmount that the salesperson gave
	// by hand, kept separately so a campaign's real cost stays measurable.
	ManualDiscount uint `bson:"manual_discount,omitempty"`

	// IsCustom marks a hand-typed invoice line (repair fee, Apple ID, a part
	// sourced from outside) with no catalog product behind it.
	IsCustom bool `bson:"is_custom,omitempty"`

	// --- inventory & accounting, captured at the moment of sale ---

	// Kind mirrors the product's kind so an invoice line stays readable even
	// if the catalog entry changes later. Empty is read as physical.
	Kind ProductKind `bson:"kind,omitempty"`
	// CostAtTime is the unit cost basis (weighted-average stock cost, or the
	// expected external supply cost) frozen when the order was placed.
	CostAtTime uint `bson:"cost_at_time,omitempty"`
	// Supply says how these units were sourced.
	Supply SupplySource `bson:"supply,omitempty"`
	// StockQuantity / ExternalQuantity split Quantity across the two sources.
	StockQuantity    uint `bson:"stock_quantity,omitempty"`
	ExternalQuantity uint `bson:"external_quantity,omitempty"`
	// SupplierID is the supplier the external part was sourced from.
	SupplierID string `bson:"supplier_id,omitempty"`
}

// GetKind normalises pre-inventory order lines.
func (i OrderItem) GetKind() ProductKind {
	if i.Kind == "" {
		return ProductKindPhysical
	}
	return i.Kind
}

// GrossProfit is revenue minus cost for this line, in the order's currency.
// Revenue is what the customer actually paid, so every discount taken off the
// line — campaign or counter — is subtracted before the cost.
// Signed: selling below cost is a real loss and must report as one.
func (i OrderItem) GrossProfit() int64 {
	revenue := int64(i.PriceAtTime)*int64(i.Quantity) - int64(i.DiscountAmount)
	return revenue - int64(i.CostAtTime)*int64(i.Quantity)
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

package orderdeliveryservice

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	orderdeliveryparam "github.com/abolfazlnorzad/steve/param/orderdelivery"
	"github.com/abolfazlnorzad/steve/pkg/phonenumber"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	ordernotifyservice "github.com/abolfazlnorzad/steve/service/ordernotify"
	qrcode "github.com/skip2/go-qrcode"
)

const confirmedByCustomerQR = "customer_qr_otp"

type OrderRepo interface {
	GetOrderByID(ctx context.Context, orderID string) (entity.Order, error)
	GetOrderByDeliveryToken(ctx context.Context, token string) (entity.Order, error)
	EnsureDeliveryToken(ctx context.Context, orderID string, token string) (string, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error
	SetDeliveryConfirmedBy(ctx context.Context, orderID string, by string) error
}

type UserRepo interface {
	GetUserByID(ctx context.Context, id string) (entity.User, error)
}

type OTPService interface {
	GenerateOTPCode(ctx context.Context, phoneNumber string, prefix string) error
	VerifyOTPCode(ctx context.Context, phoneNumber string, code string, prefix string) (bool, error)
}

type OTPStore interface {
	IncrAttempts(ctx context.Context, key string, ttl time.Duration) (int64, error)
	DeleteUserActiveCode(ctx context.Context, key string) error
	StoreNewActiveCodeForUser(ctx context.Context, key string, code string, expireTime time.Duration) error
	GetUserActiveCode(ctx context.Context, key string) (string, error)
}

type Service struct {
	cfg      Config
	orders   OrderRepo
	users    UserRepo
	otp      OTPService
	otpStore OTPStore
}

func New(cfg Config, orders OrderRepo, users UserRepo, otp OTPService, otpStore OTPStore) Service {
	return Service{cfg: cfg, orders: orders, users: users, otp: otp, otpStore: otpStore}
}

func (s Service) GetQRCode(ctx context.Context, req orderdeliveryparam.GetQRCodeRequest) (orderdeliveryparam.GetQRCodeResponse, error) {
	const op = "orderdeliveryservice.GetQRCode"

	order, err := s.orders.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		return orderdeliveryparam.GetQRCodeResponse{}, richerror.New(op).WithWarpError(err)
	}

	token := order.DeliveryToken
	if token == "" {
		fresh, gErr := newToken()
		if gErr != nil {
			return orderdeliveryparam.GetQRCodeResponse{}, richerror.New(op).WithWarpError(gErr).
				WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
		}
		if token, err = s.orders.EnsureDeliveryToken(ctx, order.ID, fresh); err != nil {
			return orderdeliveryparam.GetQRCodeResponse{}, richerror.New(op).WithWarpError(err)
		}
	}

	link := s.link(token)
	size := s.cfg.QRSize
	if size <= 0 {
		size = 512
	}
	png, qErr := qrcode.Encode(link, qrcode.Medium, size)
	if qErr != nil {
		return orderdeliveryparam.GetQRCodeResponse{}, richerror.New(op).WithWarpError(qErr).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	return orderdeliveryparam.GetQRCodeResponse{
		OrderID: order.ID,
		URL:     link,
		QRCode:  "data:image/png;base64," + base64.StdEncoding.EncodeToString(png),
	}, nil
}

func (s Service) GetInfo(ctx context.Context, req orderdeliveryparam.GetInfoRequest) (orderdeliveryparam.GetInfoResponse, error) {
	const op = "orderdeliveryservice.GetInfo"

	order, err := s.orders.GetOrderByDeliveryToken(ctx, req.Token)
	if err != nil {
		return orderdeliveryparam.GetInfoResponse{}, richerror.New(op).WithWarpError(err)
	}
	user, err := s.users.GetUserByID(ctx, order.UserID)
	if err != nil {
		return orderdeliveryparam.GetInfoResponse{}, richerror.New(op).WithWarpError(err)
	}

	return orderdeliveryparam.GetInfoResponse{
		OrderID:     order.ID,
		Status:      string(order.Status),
		StatusLabel: ordernotifyservice.StatusLabelFa(order.Status),
		MaskedPhone: maskPhone(user.PhoneNumber),
		CanConfirm:  entity.CanTransitionOrderStatus(order.Status, entity.OrderStatusDelivered),
	}, nil
}

func (s Service) SendOTP(ctx context.Context, req orderdeliveryparam.SendOTPRequest) (orderdeliveryparam.SendOTPResponse, error) {
	const op = "orderdeliveryservice.SendOTP"

	order, phone, err := s.confirmableOrder(ctx, req.Token, req.PhoneNumber)
	if err != nil {
		return orderdeliveryparam.SendOTPResponse{}, richerror.New(op).WithWarpError(err)
	}
	if err := s.otp.GenerateOTPCode(ctx, phone, s.otpPrefix(order.ID)); err != nil {
		return orderdeliveryparam.SendOTPResponse{}, richerror.New(op).WithWarpError(err)
	}
	return orderdeliveryparam.SendOTPResponse{MaskedPhone: maskPhone(phone)}, nil
}

func (s Service) Verify(ctx context.Context, req orderdeliveryparam.VerifyRequest) (orderdeliveryparam.VerifyResponse, error) {
	const op = "orderdeliveryservice.Verify"

	order, phone, err := s.confirmableOrder(ctx, req.Token, req.PhoneNumber)
	if err != nil {
		return orderdeliveryparam.VerifyResponse{}, richerror.New(op).WithWarpError(err)
	}

	prefix := s.otpPrefix(order.ID)
	attempts, aErr := s.otpStore.IncrAttempts(ctx, prefix+":attempts", s.attemptWindow())
	if aErr != nil {
		return orderdeliveryparam.VerifyResponse{}, richerror.New(op).WithWarpError(aErr)
	}
	if attempts > s.maxAttempts() {
		return orderdeliveryparam.VerifyResponse{}, richerror.New(op).
			WithKind(richerror.KindForbidden).WithMessageKey("order.delivery_too_many_attempts")
	}

	ok, vErr := s.otp.VerifyOTPCode(ctx, phone, req.Code, prefix)
	if vErr != nil {
		return orderdeliveryparam.VerifyResponse{}, richerror.New(op).WithWarpError(vErr)
	}
	if !ok {
		return orderdeliveryparam.VerifyResponse{}, richerror.New(op).
			WithKind(richerror.KindUnauthorized).WithMessageKey("validation.incorrect_opt_code")
	}
	_ = s.otpStore.DeleteUserActiveCode(ctx, fmt.Sprintf("%s:%s", prefix, phone))

	// Goes through the notifying repo, so the customer also gets the
	// "delivered" SMS.
	if err := s.orders.UpdateOrderStatus(ctx, order.ID, entity.OrderStatusDelivered); err != nil {
		return orderdeliveryparam.VerifyResponse{}, richerror.New(op).WithWarpError(err)
	}
	_ = s.orders.SetDeliveryConfirmedBy(ctx, order.ID, confirmedByCustomerQR)
	order.Status = entity.OrderStatusDelivered

	invoiceToken, gErr := newToken()
	if gErr != nil {
		return orderdeliveryparam.VerifyResponse{}, richerror.New(op).WithWarpError(gErr).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	ttl := s.invoiceTokenTTL()
	if err := s.otpStore.StoreNewActiveCodeForUser(ctx, s.invoiceKey(invoiceToken), order.ID, ttl); err != nil {
		return orderdeliveryparam.VerifyResponse{}, richerror.New(op).WithWarpError(err)
	}

	return orderdeliveryparam.VerifyResponse{
		OrderID:               order.ID,
		Status:                string(entity.OrderStatusDelivered),
		Invoice:               buildInvoice(order),
		InvoiceToken:          invoiceToken,
		InvoiceTokenExpiresIn: int64(ttl.Seconds()),
	}, nil
}

// GetInvoice returns the invoice to a customer who proved ownership with an
// OTP; the invoice token issued by Verify is bound to the QR's order.
func (s Service) GetInvoice(ctx context.Context, req orderdeliveryparam.GetInvoiceRequest) (orderdeliveryparam.GetInvoiceResponse, error) {
	const op = "orderdeliveryservice.GetInvoice"

	unauthorized := richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessageKey("order.delivery_invoice_unauthorized")
	if req.InvoiceToken == "" {
		return orderdeliveryparam.GetInvoiceResponse{}, unauthorized
	}
	order, err := s.orders.GetOrderByDeliveryToken(ctx, req.Token)
	if err != nil {
		return orderdeliveryparam.GetInvoiceResponse{}, richerror.New(op).WithWarpError(err)
	}
	orderID, err := s.otpStore.GetUserActiveCode(ctx, s.invoiceKey(req.InvoiceToken))
	if err != nil || orderID != order.ID {
		return orderdeliveryparam.GetInvoiceResponse{}, unauthorized
	}
	return orderdeliveryparam.GetInvoiceResponse{Invoice: buildInvoice(order)}, nil
}

func buildInvoice(order entity.Order) orderdeliveryparam.Invoice {
	items := make([]orderdeliveryparam.InvoiceItem, 0, len(order.Items))
	for _, it := range order.Items {
		gross := it.PriceAtTime * it.Quantity
		total := uint(0)
		if gross > it.DiscountAmount {
			total = gross - it.DiscountAmount
		}
		items = append(items, orderdeliveryparam.InvoiceItem{
			ProductName:  it.ProductName,
			VariantLabel: variantLabel(it.VariantAttributes),
			Quantity:     it.Quantity,
			UnitPrice:    it.PriceAtTime,
			Discount:     it.DiscountAmount,
			Total:        total,
		})
	}
	// Derived rather than summed so subtotal + shipping - discount always
	// equals final_amount, whatever mix of line/coupon/counter discounts applied.
	discount := uint(0)
	if gross := order.SubtotalAmount + order.ShippingAmount; gross > order.FinalAmount {
		discount = gross - order.FinalAmount
	}
	a := order.ShippingAddress
	return orderdeliveryparam.Invoice{
		OrderID:      order.ID,
		Status:       string(order.Status),
		StatusLabel:  ordernotifyservice.StatusLabelFa(order.Status),
		Items:        items,
		Subtotal:     order.SubtotalAmount,
		ShippingCost: order.ShippingAmount,
		Discount:     discount,
		FinalAmount:  order.FinalAmount,
		ShippingAddress: orderdeliveryparam.InvoiceAddress{
			Title: a.Title, Province: a.Province, City: a.City, Street: a.Street, Postcode: a.Postcode,
		},
		CreatedAt: order.CreatedAt,
	}
}

func variantLabel(attrs []entity.VariantAttribute) string {
	parts := make([]string, 0, len(attrs))
	for _, a := range attrs {
		if a.Label != "" {
			parts = append(parts, a.Label)
		} else if a.Value != "" {
			parts = append(parts, a.Value)
		}
	}
	return strings.Join(parts, " / ")
}

// confirmableOrder resolves the token, checks the order can still move to
// delivered and that phone belongs to the order's owner.
func (s Service) confirmableOrder(ctx context.Context, token, rawPhone string) (entity.Order, string, error) {
	const op = "orderdeliveryservice.confirmableOrder"

	order, err := s.orders.GetOrderByDeliveryToken(ctx, token)
	if err != nil {
		return entity.Order{}, "", richerror.New(op).WithWarpError(err)
	}
	if order.Status == entity.OrderStatusDelivered {
		return entity.Order{}, "", richerror.New(op).
			WithKind(richerror.KindBadRequest).WithMessageKey("order.delivery_already_confirmed")
	}
	if !entity.CanTransitionOrderStatus(order.Status, entity.OrderStatusDelivered) {
		return entity.Order{}, "", richerror.New(op).
			WithKind(richerror.KindBadRequest).WithMessageKey("order.delivery_not_confirmable")
	}

	user, err := s.users.GetUserByID(ctx, order.UserID)
	if err != nil {
		return entity.Order{}, "", richerror.New(op).WithWarpError(err)
	}
	phone := phonenumber.Normalize(rawPhone)
	if phone == "" || phone != phonenumber.Normalize(user.PhoneNumber) {
		return entity.Order{}, "", richerror.New(op).
			WithKind(richerror.KindForbidden).WithMessageKey("order.delivery_phone_mismatch")
	}
	return order, phone, nil
}

// link puts the token in the path (<base>/<token>); the storefront page is
// /order/delivery/:token. Old labels printed as ?token= are redirected by
// the storefront.
func (s Service) link(token string) string {
	return strings.TrimRight(s.cfg.BaseURL, "/") + "/" + url.PathEscape(token)
}

func (s Service) otpPrefix(orderID string) string {
	p := s.cfg.OTPPrefix
	if p == "" {
		p = "order_delivery"
	}
	return p + ":" + orderID
}

func (s Service) invoiceKey(invoiceToken string) string {
	p := s.cfg.OTPPrefix
	if p == "" {
		p = "order_delivery"
	}
	return p + ":invoice:" + invoiceToken
}

func (s Service) invoiceTokenTTL() time.Duration {
	if s.cfg.InvoiceTokenTTL <= 0 {
		return 30 * time.Minute
	}
	return s.cfg.InvoiceTokenTTL
}

func (s Service) maxAttempts() int64 {
	if s.cfg.MaxOTPAttempts <= 0 {
		return 5
	}
	return s.cfg.MaxOTPAttempts
}

func (s Service) attemptWindow() time.Duration {
	if s.cfg.AttemptWindow <= 0 {
		return 15 * time.Minute
	}
	return s.cfg.AttemptWindow
}

func newToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func maskPhone(p string) string {
	r := []rune(p)
	if len(r) < 7 {
		return p
	}
	return string(r[:4]) + strings.Repeat("*", len(r)-6) + string(r[len(r)-2:])
}

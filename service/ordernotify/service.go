package ordernotifyservice

import (
	"context"
	"crypto/rand"
	"math/big"
	"strconv"
	"time"

	"github.com/abolfazlnorzad/steve/adapter/mellipayamak"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type OrderRepo interface {
	ClaimThankYouCoupon(ctx context.Context, orderID string) (bool, error)
	ReleaseThankYouCoupon(ctx context.Context, orderID string) error
	SetThankYouDiscountID(ctx context.Context, orderID string, discountID string) error
}

type DiscountRepo interface {
	Create(ctx context.Context, discount entity.Discount) error
	CheckCodeIsUnique(ctx context.Context, code string) (bool, error)
}

type UserRepo interface {
	GetUserByID(ctx context.Context, id string) (entity.User, error)
}

type SMSSender interface {
	SendPatternSMS(ctx context.Context, to string, text []string, patternKey string) error
}

type Service struct {
	cfg          Config
	orders       OrderRepo
	discountRepo DiscountRepo
	userRepo     UserRepo
	sms          SMSSender
}

func New(cfg Config, orders OrderRepo, discountRepo DiscountRepo, userRepo UserRepo, sms SMSSender) Service {
	return Service{cfg: cfg, orders: orders, discountRepo: discountRepo, userRepo: userRepo, sms: sms}
}

// StatusChanged is called after an order's status was persisted. It runs
// detached from the request: an SMS failure must never fail the status
// change itself.
func (s Service) StatusChanged(order entity.Order, to entity.OrderStatus) {
	if !s.cfg.Enabled || order.Status == to || order.UserID == "" {
		return
	}
	// An unpaid order expiring is not news to the customer.
	if order.Status == entity.OrderStatusPendingPayment && to == entity.OrderStatusCancelled {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.smsTimeout())
		defer cancel()

		defer func() {
			if r := recover(); r != nil {
				logger.L().Error("ordernotify_panic", "order_id", order.ID, "panic", r)
			}
		}()

		user, err := s.userRepo.GetUserByID(ctx, order.UserID)
		if err != nil || user.PhoneNumber == "" {
			logger.L().Error("ordernotify_user_lookup_failed", "order_id", order.ID, "err", err)
			return
		}

		placed := order.Status == entity.OrderStatusPendingPayment &&
			(to == entity.OrderStatusProcessing || to == entity.OrderStatusDelivered)
		if placed && !order.IsPartnerOrder && s.sendThankYou(ctx, order, user) {
			return
		}
		if to == entity.OrderStatusShipped && order.TrackingCode != "" && s.sendShipped(ctx, order, user) {
			return
		}
		s.sendStatus(ctx, order, user, to)
	}()
}

// sendThankYou mints the one-time coupon and sends the thank-you SMS.
// It returns false when nothing was sent, so the caller falls back to the
// plain status message.
func (s Service) sendThankYou(ctx context.Context, order entity.Order, user entity.User) bool {
	const op = "ordernotifyservice.sendThankYou"

	if !s.cfg.ThankYouDiscountEnabled || s.cfg.ThankYouDiscountValue == 0 {
		return false
	}

	claimed, err := s.orders.ClaimThankYouCoupon(ctx, order.ID)
	if err != nil || !claimed {
		if err != nil {
			logger.L().Error(op, "order_id", order.ID, "err", err)
		}
		return false
	}

	d, err := s.mintCoupon(ctx, user.ID)
	if err != nil {
		logger.L().Error(op, "order_id", order.ID, "err", err)
		_ = s.orders.ReleaseThankYouCoupon(ctx, order.ID)
		return false
	}
	if sErr := s.orders.SetThankYouDiscountID(ctx, order.ID, d.ID); sErr != nil {
		logger.L().Error(op, "order_id", order.ID, "err", sErr)
	}

	if sErr := s.sms.SendPatternSMS(ctx, user.PhoneNumber,
		[]string{displayName(user), order.ID, d.Code}, mellipayamak.OrderThanksKey); sErr != nil {
		logger.L().Error(op, "order_id", order.ID, "err", sErr)
	}
	return true
}

// sendShipped sends the tracking code of a shipped order. It returns false
// when the SMS could not be sent (e.g. the order_shipped pattern is not
// configured), so the caller falls back to the plain status message.
func (s Service) sendShipped(ctx context.Context, order entity.Order, user entity.User) bool {
	const op = "ordernotifyservice.sendShipped"

	if err := s.sms.SendPatternSMS(ctx, user.PhoneNumber,
		[]string{displayName(user), order.ID, order.TrackingCode}, mellipayamak.OrderShippedKey); err != nil {
		logger.L().Error(op, "order_id", order.ID, "err", err)
		return false
	}
	return true
}

func (s Service) sendStatus(ctx context.Context, order entity.Order, user entity.User, to entity.OrderStatus) {
	const op = "ordernotifyservice.sendStatus"

	if err := s.sms.SendPatternSMS(ctx, user.PhoneNumber,
		[]string{displayName(user), order.ID, StatusLabelFa(to)}, mellipayamak.OrderStatusKey); err != nil {
		logger.L().Error(op, "order_id", order.ID, "status", to, "err", err)
	}
}

func (s Service) mintCoupon(ctx context.Context, userID string) (entity.Discount, error) {
	const op = "ordernotifyservice.mintCoupon"

	dtype := entity.DiscountType(s.cfg.ThankYouDiscountType)
	if dtype != entity.PercentageDiscount && dtype != entity.FixedAmountDiscount {
		dtype = entity.FixedAmountDiscount
	}

	var code string
	for attempt := 0; attempt < 6; attempt++ {
		c, gErr := generateCouponCode()
		if gErr != nil {
			return entity.Discount{}, richerror.New(op).WithWarpError(gErr)
		}
		unique, uErr := s.discountRepo.CheckCodeIsUnique(ctx, c)
		if uErr != nil {
			return entity.Discount{}, richerror.New(op).WithWarpError(uErr)
		}
		if unique {
			code = c
			break
		}
	}
	if code == "" {
		return entity.Discount{}, richerror.New(op).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	d := entity.Discount{
		Code:                  code,
		Type:                  dtype,
		Value:                 s.cfg.ThankYouDiscountValue,
		MinimumPurchaseAmount: s.cfg.ThankYouMinimumPurchaseAmount,
		StartsAt:              time.Now().Add(-time.Minute),
		ExpiresAt:             time.Now().AddDate(0, 0, s.cfg.validDays()),
		TotalUsageLimit:       1,
		UsageLimitPerUser:     1,
		ApplicableTo:          entity.ApplicableToUsers,
		ApplicableUserIDs:     []string{userID},
		IsActive:              true,
	}
	d.BeforeCreate()
	if err := s.discountRepo.Create(ctx, d); err != nil {
		return entity.Discount{}, richerror.New(op).WithWarpError(err)
	}
	return d, nil
}

const codeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func generateCouponCode() (string, error) {
	const n = 8
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeAlphabet))))
		if err != nil {
			return "", err
		}
		b[i] = codeAlphabet[idx.Int64()]
	}
	return "THX" + string(b), nil
}

func displayName(u entity.User) string {
	if u.Name != "" {
		return u.Name
	}
	return "مشتری"
}

// StatusLabelFa is the customer-facing Persian name of a status.
func StatusLabelFa(st entity.OrderStatus) string {
	switch st {
	case entity.OrderStatusPendingPayment:
		return "در انتظار پرداخت"
	case entity.OrderStatusProcessing:
		return "در حال پردازش"
	case entity.OrderStatusShipped:
		return "ارسال شده"
	case entity.OrderStatusDelivered:
		return "تحویل داده شده"
	case entity.OrderStatusCancelled:
		return "لغو شده"
	case entity.OrderStatusRefunded:
		return "مرجوع شده"
	default:
		return strconv.Quote(string(st))
	}
}

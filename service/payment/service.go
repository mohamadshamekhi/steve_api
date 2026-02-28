package paymentservice

import (
	"context"
	"fmt"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	paymentparam "github.com/abolfazlnorzad/steve/param/payment"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/go-telegram/bot"
	persian "github.com/yaa110/go-persian-calendar"
)

type Config struct {
	CallbackURL                             string        `koanf:"callback_url"`
	Gateway                                 string        `koanf:"gateway"`
	ExpirationTimePaymentSmsContext         time.Duration `koanf:"expiration_time_payment_sms_context"`
	ExpirationTimePaymentMsgTelegramContext time.Duration `koanf:"expiration_time_payment_msg_telegram_context"`
	RetryCount                              int           `koanf:"retry_count"`
	RetryDelayTime                          time.Duration `koanf:"retry_delay_time"`
}
type GatewayAdapter interface {
	RequestPayment(amount uint64, callbackURL string, description string) (authority string, paymentURL string, err error)
	VerifyPayment(amount uint64, authority string) (string, string, error)
}

type TelegramBotAdapter interface {
	SendMessageAdmin(ctx context.Context, text string) error
	SendMessageDeveloper(ctx context.Context, text string) error
	GetBot() *bot.Bot
}

type Repository interface {
	CreatePayment(ctx context.Context, p entity.Payment) error
	UpdateGatewayTransactionID(ctx context.Context, paymentID string, transactionID string) error
	GetPaymentByGatewayTransactionID(ctx context.Context, authority string) (entity.Payment, error)
	UpdatePaymentSuccess(ctx context.Context, paymentID string, refID string, cardPan string) error
	UpdatePaymentFailed(ctx context.Context, paymentID string) error
	GetPaymentsByOrderID(ctx context.Context, orderID string) ([]entity.Payment, error)
}

type DiscountProvider interface {
	UpdateDiscountUsageAfterRegisterOrder(ctx context.Context, req discountparam.UpdateDiscountUsageAfterRegisterOrderRequest) (discountparam.UpdateDiscountUsageAfterRegisterOrderResponse, error)
}

type OrderProvider interface {
	GetOrderByID(ctx context.Context, req orderparam.GetOrderByIDRequest) (orderparam.GetOrderByIDResponse, error)
	UpdateOrderStatus(ctx context.Context, request orderparam.UpdateStatusRequest) (orderparam.UpdateStatusResponse, error)
}

type ProductProvider interface {
	IncrementSalesCount(ctx context.Context, req productparam.UpdateProductSalesCountRequest) (productparam.UpdateProductSalesCountResponse, error)
}

type MelliPayamakProvider interface {
	SendPayment(ctx context.Context, to string, text []string) error
}

type UserProvider interface {
	GetUserByID(ctx context.Context, req userparam.GetUsersRequest) (userparam.GetUsersResponse, error)
}

type Service struct {
	cfg                         Config
	repo                        Repository
	orderProvider               OrderProvider
	gatewayAdapter              GatewayAdapter
	discountProvider            DiscountProvider
	productProvider             ProductProvider
	melliPayamakProvider        MelliPayamakProvider
	userProvider                UserProvider
	telegramAdminBotAdapter     TelegramBotAdapter
	telegramDeveloperBotAdapter TelegramBotAdapter
}

func NewService(cfg Config, repo Repository,
	orderProvider OrderProvider,
	gatewayAdapter GatewayAdapter,
	discountProvider DiscountProvider,
	productProvider ProductProvider,
	melliPayamakProvider MelliPayamakProvider,
	userProvider UserProvider,
	telegramAdminBotAdapter TelegramBotAdapter,
	telegramDeveloperBotAdapter TelegramBotAdapter,
) Service {
	return Service{
		cfg:                         cfg,
		repo:                        repo,
		orderProvider:               orderProvider,
		gatewayAdapter:              gatewayAdapter,
		discountProvider:            discountProvider,
		productProvider:             productProvider,
		melliPayamakProvider:        melliPayamakProvider,
		userProvider:                userProvider,
		telegramAdminBotAdapter:     telegramAdminBotAdapter,
		telegramDeveloperBotAdapter: telegramDeveloperBotAdapter,
	}
}

func (s Service) GoToGateway(ctx context.Context, req paymentparam.GoToGatewayRequest) (paymentparam.GoToGatewayResponse, error) {
	const op = "paymentservice.GoToGateway"
	order, oErr := s.orderProvider.GetOrderByID(ctx, orderparam.GetOrderByIDRequest{ID: req.OrderID})
	if oErr != nil {
		return paymentparam.GoToGatewayResponse{}, richerror.New(op).WithWarpError(oErr)
	}
	if order.Order.UserID != req.UserID {
		return paymentparam.GoToGatewayResponse{}, richerror.New(op).
			WithKind(richerror.KindForbidden).
			WithMessageKey("payment.order_does_not_belong_to_user")
	}
	if order.Order.Status != entity.OrderStatusPendingPayment {
		return paymentparam.GoToGatewayResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("payment.order_already_paid_or_cancelled")
	}
	p := entity.Payment{
		UserID:  req.UserID,
		OrderID: req.OrderID,
		Gateway: s.cfg.Gateway,
		Amount:  uint64(order.Order.FinalAmount),
	}
	p.BeforeCreate()
	err := s.repo.CreatePayment(ctx, p)
	if err != nil {
		return paymentparam.GoToGatewayResponse{}, richerror.New(op).WithWarpError(err)
	}

	tID, pLink, err := s.gatewayAdapter.RequestPayment(p.Amount, s.cfg.CallbackURL, p.Description)
	if err != nil {
		return paymentparam.GoToGatewayResponse{}, richerror.New(op).WithWarpError(err)
	}

	if err := s.repo.UpdateGatewayTransactionID(ctx, p.ID, tID); err != nil {
		return paymentparam.GoToGatewayResponse{}, richerror.New(op).WithWarpError(err)
	}

	// update order status ??

	return paymentparam.GoToGatewayResponse{
		GatewayTransactionID: tID,
		PaymentURL:           pLink,
	}, nil
}

func (s Service) VerifyPayment(ctx context.Context, req paymentparam.VerifyPaymentRequest) (paymentparam.VerifyPaymentResponse, error) {
	const op = "paymentservice.VerifyPayment"
	payment, err := s.repo.GetPaymentByGatewayTransactionID(ctx, req.GatewayTransactionID)
	if err != nil {
		return paymentparam.VerifyPaymentResponse{}, richerror.New(op).WithWarpError(err)
	}

	if payment.Status == entity.PaymentStatusSuccessful {
		return paymentparam.VerifyPaymentResponse{
			Status:       entity.PaymentStatusSuccessful,
			TrackingCode: payment.TrackingCode,
			Message:      "تراکنش قبلاً با موفقیت انجام شده است",
		}, nil
	}
	refID, cardPan, err := s.gatewayAdapter.VerifyPayment(payment.Amount, req.GatewayTransactionID)
	if err != nil {
		_ = s.repo.UpdatePaymentFailed(ctx, payment.ID)
		logger.L().Error("important , gateway", err)
		return paymentparam.VerifyPaymentResponse{}, richerror.New(op).WithWarpError(err)
	}
	if err := s.repo.UpdatePaymentSuccess(ctx, payment.ID, refID, cardPan); err != nil {
		logger.L().Error("CRITICAL: Payment verified but DB update failed", "payment_id", payment.ID, "err", err)
	}
	order, err := s.orderProvider.GetOrderByID(ctx, orderparam.GetOrderByIDRequest{ID: payment.OrderID})
	if err != nil {
		logger.L().Error("CRITICAL: Failed to get order after payment", "order_id", payment.OrderID, "err", err)
	} else {
		if _, err := s.orderProvider.UpdateOrderStatus(ctx, orderparam.UpdateStatusRequest{
			OrderID: payment.OrderID,
			Status:  entity.OrderStatusProcessing,
		}); err != nil {
			logger.L().Error("Failed to update order status", "order_id", order.Order.ID, "err", err)
		}

		if len(order.Order.AppliedDiscounts) > 0 {
			_, err = s.discountProvider.UpdateDiscountUsageAfterRegisterOrder(ctx, discountparam.UpdateDiscountUsageAfterRegisterOrderRequest{
				AppliedDiscount: order.Order.AppliedDiscounts,
			})
			if err != nil {
				logger.L().Error("Failed to update discount usage", "order_id", order.Order.ID, "err", err)
			}
		}

		for _, item := range order.Order.Items {
			if _, err = s.productProvider.IncrementSalesCount(ctx, productparam.UpdateProductSalesCountRequest{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			}); err != nil {
				logger.L().Error("Failed to increment product sales count", "product_id", item.ProductID, "err", err)
			}
		}
	}

	userRes, err := s.userProvider.GetUserByID(ctx, userparam.GetUsersRequest{ID: payment.UserID})
	if err != nil {
		logger.L().Error("Failed to get user", "user_id", payment.UserID, "err", err)
	}

	loc, _ := time.LoadLocation("Asia/Tehran")
	now := time.Now().In(loc)
	pt := persian.New(now)
	date := pt.Format("yyyy/MM/dd")
	hour := pt.Format("HH:mm")

	go func(date string, hour string) {
		otpCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ExpirationTimePaymentSmsContext)
		defer cancel()

		s.retry(
			otpCtx,
			map[string]any{
				"op":          "send_sms_payment",
				"username":    userRes.User.Name,
				"phoneNumber": userRes.User.PhoneNumber,
				"refID":       refID,
				"date":        date + " " + hour,
			},
			func() error {
				return s.melliPayamakProvider.SendPayment(otpCtx, userRes.User.PhoneNumber, []string{userRes.User.Name, date, hour, refID})
			},
		)
	}(date, hour)

	go func(date string, hour string) {
		tCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ExpirationTimePaymentMsgTelegramContext)
		defer cancel()

		amount := payment.Amount
		name := userRes.User.Name
		phoneNumber := userRes.User.PhoneNumber
		shippingName := order.Order.Shipping.Name
		shippingAddress := fmt.Sprintf("%s %s %s", order.Order.ShippingAddress.Province, order.Order.ShippingAddress.City, order.Order.ShippingAddress.Street)
		orderID := order.Order.ID

		var itemsText string
		for i, item := range order.Order.Items {
			itemsText += fmt.Sprintf(
				"\n%d- نام محصول: %s\nتعداد محصول: %d\n",
				i+1,
				item.ProductName,
				item.Quantity,
			)
		}
		message := fmt.Sprintf(`✅ یک سفارش با موفقیت ثبت شد

		🔹 کد پیگیری: %s
		🔹 مبلغ خرید: %v
		🔹 نام کاربر: %s
		🔹 شماره موبایل: %s
		🔹 تاریخ ثبت: %s
		🔹 ساعت ثبت: %s
		🔹 روش ارسال: %s
		🔹 آدرس: %s

		📦 لیست آیتم‌ها:
		%s
		https://backoffice.stevemobile.com/dashboard/order?order=%s
		`, refID, amount, name, phoneNumber, date, hour, shippingName, shippingAddress, itemsText, orderID)

		if s.telegramAdminBotAdapter.GetBot() != nil {
			s.retry(
				tCtx,
				map[string]any{
					"op":        "send_msg_telegram_admin_bot",
					"refID":     refID,
					"amount":    amount,
					"name":      name,
					"phone":     phoneNumber,
					"date":      date + " " + hour,
					"orderID":   orderID,
					"paymentID": payment.ID,
				},
				func() error {
					return s.telegramAdminBotAdapter.SendMessageAdmin(tCtx, message)
				},
			)
		} else {
			logger.L().Error("Telegram bot failed to sent message to admin chat")
		}

	}(date, hour)

	return paymentparam.VerifyPaymentResponse{
		Status:       entity.PaymentStatusSuccessful,
		TrackingCode: refID,
		Message:      "پرداخت با موفقیت انجام شد",
	}, nil
}

func (s Service) GetPaymentsByOrderID(ctx context.Context, req paymentparam.GetPaymentsByOrderIDRequest) (paymentparam.GetPaymentsByOrderIDResponse, error) {
	const op = "paymentservice.GetPaymentsByOrderID"

	payments, err := s.repo.GetPaymentsByOrderID(ctx, req.OrderID)
	if err != nil {
		return paymentparam.GetPaymentsByOrderIDResponse{}, richerror.New(op).WithWarpError(err)
	}

	return paymentparam.GetPaymentsByOrderIDResponse{Payments: payments}, nil
}

func (s Service) retry(ctx context.Context, meta map[string]any, fn func() error,
) {
	var lastErr error

	for i := 0; i < s.cfg.RetryCount; i++ {
		err := fn()
		if err == nil {
			return
		}

		lastErr = err

		logger.L().Error("retry failed",
			"attempt", i+1,
			"err", err,
		)

		if i == s.cfg.RetryCount-1 {
			go s.notifyDeveloper(lastErr, meta)
			return
		}

		select {
		case <-time.After(s.cfg.RetryDelayTime):
		case <-ctx.Done():
			return
		}
	}
}

func (s Service) notifyDeveloper(err error, meta map[string]any) {
	tCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ExpirationTimePaymentMsgTelegramContext)
	defer cancel()

	var metaText string
	for k, v := range meta {
		metaText += fmt.Sprintf(
			"\n🔹 %s: %v\n",
			k,
			v,
		)
	}

	message := fmt.Sprintf(`
	🚨 One Bug Found :
	%s

	Error : %v
	`, metaText, err)

	if s.telegramDeveloperBotAdapter.GetBot() != nil {
		if err := s.telegramDeveloperBotAdapter.SendMessageDeveloper(tCtx, message); err != nil {
			logger.L().Error("Telegram bot failed to send message to developer chat", "err: ", err)
		}
	}

func (s Service) CreatePaymentForBackofficePanel(ctx context.Context, req paymentparam.CreatePaymentForBackofficePanelRequest) (paymentparam.CreatePaymentForBackofficePanelResponse, error) {
	const op = "paymentservice.CreatePaymentForBackofficePanel"
	order, err := s.orderProvider.GetOrderByID(ctx, orderparam.GetOrderByIDRequest{ID: req.OrderID})
	if err != nil {
		return paymentparam.CreatePaymentForBackofficePanelResponse{}, richerror.New(op).WithWarpError(err)
	}
	payment := entity.Payment{
		UserID:      order.Order.UserID,
		OrderID:     order.Order.ID,
		Amount:      uint64(order.Order.FinalAmount),
		Gateway:     "in_person",
		Description: req.Description,
	}
	payment.BeforeCreate()
	payment.Status = entity.PaymentStatusSuccessful

	pErr := s.repo.CreatePayment(ctx, payment)
	if pErr != nil {
		return paymentparam.CreatePaymentForBackofficePanelResponse{}, richerror.New(op).WithWarpError(pErr)
	}
	return paymentparam.CreatePaymentForBackofficePanelResponse{}, nil

}

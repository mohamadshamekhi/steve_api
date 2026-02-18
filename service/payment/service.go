package paymentservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	paymentparam "github.com/abolfazlnorzad/steve/param/payment"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	persian "github.com/yaa110/go-persian-calendar"
)

type Config struct {
	CallbackURL                     string        `koanf:"callback_url"`
	Gateway                         string        `koanf:"gateway"`
	ExpirationTimePaymentSmsContext time.Duration `koanf:"expiration_time_payment_sms_context"`
}
type GatewayAdapter interface {
	RequestPayment(amount uint64, callbackURL string, description string) (authority string, paymentURL string, err error)
	VerifyPayment(amount uint64, authority string) (string, string, error)
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
	cfg                  Config
	repo                 Repository
	orderProvider        OrderProvider
	gatewayAdapter       GatewayAdapter
	discountProvider     DiscountProvider
	productProvider      ProductProvider
	melliPayamakProvider MelliPayamakProvider
	userProvider         UserProvider
}

func NewService(cfg Config, repo Repository, orderProvider OrderProvider, gatewayAdapter GatewayAdapter, discountProvider DiscountProvider, productProvider ProductProvider, melliPayamakProvider MelliPayamakProvider, userProvider UserProvider) Service {
	return Service{
		cfg:                  cfg,
		repo:                 repo,
		orderProvider:        orderProvider,
		gatewayAdapter:       gatewayAdapter,
		discountProvider:     discountProvider,
		productProvider:      productProvider,
		melliPayamakProvider: melliPayamakProvider,
		userProvider:         userProvider,
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

	user, err := s.userProvider.GetUserByID(ctx, userparam.GetUsersRequest{ID: payment.UserID})
	if err != nil {
		logger.L().Error("Failed to get user", "user_id", payment.UserID, "err", err)
	}

	go func() {
		otpCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ExpirationTimePaymentSmsContext)
		defer cancel()
		loc, _ := time.LoadLocation("Asia/Tehran")
		now := time.Now().In(loc)
		pt := persian.New(now)
		date := pt.Format("yyyy/MM/dd")
		hour := pt.Format("HH:mm")

		if err := s.melliPayamakProvider.SendPayment(otpCtx, user.User.PhoneNumber, []string{user.User.Name, date, hour, refID}); err != nil {
			logger.L().Error("Failed to send payment SMS", "err", err)
		}
	}()

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

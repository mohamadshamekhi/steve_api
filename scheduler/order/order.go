package orderscheduler

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/robfig/cron/v3"
)

type OrderScheduler struct {
	cron          *cron.Cron
	cfg           Config
	OrderProvider OrderSvcProvider
}

type Config struct {
	CanselExpiredPendingOrdersTime time.Duration `koanf:"cancel_expired_pending_orders_time"`
	OrderSchedulerCycleTime        time.Duration `koanf:"order_scheduler_cycle_time"`
}

type OrderSvcProvider interface {
	CancelExpiredPendingOrders(ctx context.Context, req orderparam.CancelExpiredPendingOrdersRequest) (orderparam.CancelExpiredPendingOrdersResponse, error)
}

func New(cfg Config, orderProvider OrderSvcProvider) *OrderScheduler {
	return &OrderScheduler{
		cron:          cron.New(),
		cfg:           cfg,
		OrderProvider: orderProvider,
	}
}

func (s *OrderScheduler) OrderSchedulerRun() {
	_, err := s.cron.AddFunc(
		"@every "+s.cfg.OrderSchedulerCycleTime.String(),
		s.cancelExpiredOrders,
	)
	if err != nil {
		logger.L().Error("failed to register order scheduler", "error", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	s.cron.Start()

	<-sig

	fmt.Println("Stopping Order Scheduler...")

	stopCtx := s.cron.Stop()
	<-stopCtx.Done()
	fmt.Println("Order Scheduler stopped cleanly")
}

func (s OrderScheduler) cancelExpiredOrders() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.OrderProvider.CancelExpiredPendingOrders(ctx, orderparam.CancelExpiredPendingOrdersRequest{CanselExpiredPendingOrders: s.cfg.CanselExpiredPendingOrdersTime})
	if err != nil {
		logger.L().Error("failed to cancel expired pending orders", "error", err)
	}
}

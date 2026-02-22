package cartscheduler

import (
	"context"
	"fmt"
	"time"

	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/robfig/cron/v3"
)

type CartScheduler struct {
	cron         *cron.Cron
	cfg          Config
	cartProvider CartProvider
}

type Config struct {
	CanselExpiredPendingCartsTime time.Duration `koanf:"cancel_expired_pending_carts_time"`
	CartSchedulerCycleTime        time.Duration `koanf:"cart_scheduler_cycle_time"`
}

type CartProvider interface {
	CancelExpiredPendingCart(ctx context.Context, req cartparam.CancelExpiredPendingCartRequest) (cartparam.CancelExpiredPendingCartResponse, error)
}

func New(cfg Config, cartProvider CartProvider) *CartScheduler {
	return &CartScheduler{
		cron:         cron.New(),
		cfg:          cfg,
		cartProvider: cartProvider,
	}
}

func (s *CartScheduler) CartSchedulerRun(ctx context.Context) {
	_, err := s.cron.AddFunc(
		"@every "+s.cfg.CartSchedulerCycleTime.String(),
		s.canselExpiredCarts,
	)
	if err != nil {
		logger.L().Error("failed to add cart scheduler", "error", err)
	}

	s.cron.Start()

	<-ctx.Done()

	fmt.Println("Stopping Cart Scheduler...")

	stopCtx := s.cron.Stop()
	<-stopCtx.Done()
	fmt.Println("Cart Scheduler stopped cleanly")
}

func (s *CartScheduler) canselExpiredCarts() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.cartProvider.CancelExpiredPendingCart(ctx, cartparam.CancelExpiredPendingCartRequest{
		CanselExpiredPendingCart: s.cfg.CanselExpiredPendingCartsTime})
	if err != nil {
		logger.L().Error("failed to cancel expired pending carts", "error", err)
	}
}

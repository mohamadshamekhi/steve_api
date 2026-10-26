package birthdayscheduler

import (
	"context"
	"fmt"
	"time"

	birthdayparam "github.com/abolfazlnorzad/steve/param/birthday"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/robfig/cron/v3"
)

// tehranTZ is fixed at +03:30; Iran has no DST since 2022.
var tehranTZ = time.FixedZone("IRST", 3*3600+30*60)

type Config struct {
	// RunAt is a standard 5-field cron spec evaluated in Tehran time,
	// e.g. "0 10 * * *" = every day at 10:00.
	RunAt string `koanf:"run_at"`
}

type BirthdaySvcProvider interface {
	Run(ctx context.Context) (birthdayparam.RunResponse, error)
}

type BirthdayScheduler struct {
	cron *cron.Cron
	cfg  Config
	svc  BirthdaySvcProvider
}

func New(cfg Config, svc BirthdaySvcProvider) *BirthdayScheduler {
	return &BirthdayScheduler{
		cron: cron.New(cron.WithLocation(tehranTZ)),
		cfg:  cfg,
		svc:  svc,
	}
}

func (s *BirthdayScheduler) BirthdaySchedulerRun(ctx context.Context) {
	if _, err := s.cron.AddFunc(s.cfg.RunAt, s.run); err != nil {
		logger.L().Error("failed to register birthday scheduler", "error", err, "run_at", s.cfg.RunAt)
	}

	s.cron.Start()

	<-ctx.Done()

	fmt.Println("Stopping Birthday Scheduler...")

	stopCtx := s.cron.Stop()
	<-stopCtx.Done()
	fmt.Println("Birthday Scheduler stopped cleanly")
}

func (s *BirthdayScheduler) run() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	resp, err := s.svc.Run(ctx)
	if err != nil {
		logger.L().Error("birthday scheduler run failed", "error", err)
		return
	}
	logger.L().Info("birthday scheduler run finished", "enabled", resp.Enabled,
		"matched", resp.Matched, "granted", resp.Granted, "skipped", resp.Skipped, "failed", resp.Failed)
}

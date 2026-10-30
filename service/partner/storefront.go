package partnerservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	partnerparam "github.com/abolfazlnorzad/steve/param/partner"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	adminnotifyservice "github.com/abolfazlnorzad/steve/service/adminnotify"
)

// StorefrontService is the application entry point: a normal storefront user
// submits a partnership request and checks its status. After approval the
// partner switches to the backoffice panel (see ShopService).
type StorefrontService struct {
	repo     Repository
	cfg      Config
	telegram TelegramNotifier
	notifier AdminNotifier
}

// AdminNotifier tells the shop's admins about a customer submission by SMS.
type AdminNotifier interface {
	Notify(e adminnotifyservice.Event)
}

// SetAdminNotifier wires the admin SMS in after construction.
func (s *StorefrontService) SetAdminNotifier(n AdminNotifier) {
	s.notifier = n
}

func NewStorefrontService(repo Repository, cfg Config, telegram TelegramNotifier) StorefrontService {
	return StorefrontService{repo: repo, cfg: cfg, telegram: telegram}
}

func (s StorefrontService) programEnabled(op string) error {
	if !s.cfg.Enabled {
		return richerror.New(op).WithKind(richerror.KindForbidden).WithMessageKey("partner.program_disabled")
	}
	return nil
}

// Apply records a new partnership application, or re-opens a previously rejected
// one. A user with a pending or approved record cannot apply again.
func (s StorefrontService) Apply(ctx context.Context, userID string, req partnerparam.ApplyRequest) (partnerparam.ApplyResponse, error) {
	const op = "partnerservice.Apply"
	if err := s.programEnabled(op); err != nil {
		return partnerparam.ApplyResponse{}, err
	}

	existing, err := s.repo.FindByUserID(ctx, userID)
	if err != nil && richerror.GetKind(err) != richerror.KindNotFound {
		return partnerparam.ApplyResponse{}, richerror.New(op).WithWarpError(err)
	}

	if err == nil {
		switch existing.Status {
		case entity.PartnerStatusPending, entity.PartnerStatusApproved, entity.PartnerStatusSuspended:
			return partnerparam.ApplyResponse{}, richerror.New(op).
				WithKind(richerror.KindBadRequest).
				WithMessageKey("partner.already_applied")
		case entity.PartnerStatusRejected:
			existing.Status = entity.PartnerStatusPending
			existing.CompanyName = req.CompanyName
			existing.NationalID = req.NationalID
			existing.BusinessField = req.BusinessField
			existing.City = req.City
			existing.RejectionReason = ""
			existing.ReviewedBy = ""
			existing.ReviewedAt = nil
			existing.UpdatedAt = time.Now()
			if uErr := s.repo.Update(ctx, existing); uErr != nil {
				return partnerparam.ApplyResponse{}, richerror.New(op).WithWarpError(uErr)
			}
			s.notifyAdminNewApplication(existing)
			return partnerparam.ApplyResponse{ID: existing.ID, Status: string(existing.Status)}, nil
		}
	}

	p := entity.Partner{
		UserID:        userID,
		CompanyName:   req.CompanyName,
		NationalID:    req.NationalID,
		BusinessField: req.BusinessField,
		City:          req.City,
	}
	p.BeforeCreate()
	if cErr := s.repo.Create(ctx, p); cErr != nil {
		return partnerparam.ApplyResponse{}, richerror.New(op).WithWarpError(cErr)
	}
	s.notifyAdminNewApplication(p)

	return partnerparam.ApplyResponse{ID: p.ID, Status: string(p.Status)}, nil
}

func (s StorefrontService) notifyAdminNewApplication(p entity.Partner) {
	if s.notifier != nil {
		s.notifier.Notify(adminnotifyservice.Event{
			Title:  "درخواست همکاری",
			UserID: p.UserID,
			Detail: strings.TrimSpace(p.CompanyName + " - " + p.City),
		})
	}
	if s.telegram == nil {
		return
	}
	go func() {
		nCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		msg := fmt.Sprintf("🤝 درخواست همکاری جدید\nشرکت: %s\nشهر: %s\nحوزه: %s", p.CompanyName, p.City, p.BusinessField)
		if err := s.telegram.SendMessageAdmin(nCtx, msg); err != nil {
			logger.L().Error("partner: failed to notify admin about new application", "err", err)
		}
	}()
}

// Me returns the caller's partnership state for the storefront (application) side.
func (s StorefrontService) Me(ctx context.Context, userID string) (partnerparam.MeResponse, error) {
	const op = "partnerservice.Me"
	p, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		if richerror.GetKind(err) == richerror.KindNotFound {
			return partnerparam.MeResponse{HasPartnership: false}, nil
		}
		return partnerparam.MeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return partnerparam.MeResponse{
		HasPartnership:  true,
		Status:          string(p.Status),
		DiscountPercent: p.DiscountPercent,
		RejectionReason: p.RejectionReason,
		CompanyName:     p.CompanyName,
	}, nil
}
